package healer

import (
	"errors"
	"sort"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	log "github.com/Sirupsen/logrus"
)


const (
	MaxUint = ^uint(0)
 	MinEvacuationRangeValue = int(MaxUint >> 1) // MaxInt will be min priority
)


var FlavorsCache map[string] flavors.Flavor

func init() {
	FlavorsCache = make(map[string] flavors.Flavor)
}

type HypervisorFreeResources struct {
	HypervisorHostname string		`mapstructure:"hypervisor_hostname"`
	Id int
	State string
	Status string
	// Vcpus - VcpuUsed
	FreeVcpus int16					`mapstructure:"free_vcpus"`
	HostIP string 					`mapstructure:"host_ip"`
	FreeDiskGB int16				`mapstructure:"free_disk_gb"`
	// disk_available_least = free_disk_gb - disk_overcommit_size
	// disk_overcommit_size = virtual size of disks of all instance - used disk
	// size of all instances
	DistAvailableLeast int16		`mapstructure:"disk_available_least"`
	FreeRamMB int32					`mapstructure:"free_ram_mb"`
}

type Cluster struct {
	Resources map[string]HypervisorFreeResources
}

// Creates a slice of every host resources
func (c *Cluster) UpdateAvailableClusterResources(client *gophercloud.ServiceClient) error {
	list, err := hypervisors.GetDetailsList(client).ExtractDetails()
	if err != nil {
		return err
	}
	for _, hypervisorDetails := range list {
		res := HypervisorFreeResources{}
		res.HypervisorHostname = hypervisorDetails.HypervisorHostname
		res.Id = hypervisorDetails.Id
		res.State = hypervisorDetails.State
		res.Status = hypervisorDetails.Status
		res.FreeVcpus = hypervisorDetails.Vcpus - hypervisorDetails.VcpuUsed
		res.HostIP = hypervisorDetails.HostIP
		res.FreeDiskGB = hypervisorDetails.FreeDiskGB
		res.DistAvailableLeast = hypervisorDetails.DistAvailableLeast
		res.FreeRamMB = hypervisorDetails.FreeRamMB
		c.Resources = make(map[string]HypervisorFreeResources)
		c.Resources[hypervisorDetails.HypervisorHostname] = res
	}
	return nil
}

type ResourcesClaim struct {
	// UUID of VM
	ServerUUID string
	FlavorId   string		`mapstructure: flavor_id`
	Vcpus      int
	DiskGB     int			`mapstructure:"disk_gb"`
	RamMB      int			`mapstructure:"ram_mb"`
	RXTXFactor float64		`mapstructure: rxtx_factor`
}

func NewResourcesClaim(server servers.Server) (*ResourcesClaim, error) {
	flavorId, ok := server.Flavor["id"].(string)
	if !ok {
		return nil, errors.New("Could not create the claim. Reason flavor id convertion to string wasn't successful.")
	}
	return &ResourcesClaim{
		ServerUUID: server.ID,
		FlavorId:	flavorId,
	}, nil
}

// Creates claim of resources for VM booting
func (r* ResourcesClaim) ClaimResources(client *gophercloud.ServiceClient) error{
	flavor, ok:= FlavorsCache[r.FlavorId]
	if !ok {
		flavorPointer, err := flavors.Get(client, r.FlavorId).Extract()
		if err != nil {
			return err
		}
		flavor = *flavorPointer
	}

	r.Vcpus = flavor.VCPUs
	r.DiskGB = flavor.Disk
	r.RamMB = flavor.RAM
	r.RXTXFactor = flavor.RxTxFactor

	FlavorsCache[flavor.ID] = flavor

	return nil
}

// Helper to manage resource claims
type ResourcesClaimManager struct {
	// Mapping of instance ID and resources it needs
	ResourcesClaims 	map[string]ResourcesClaim
	TotallyUsedVcpus	int
	TotallyUsedDiskGB	int
	TotallyUsedRamMB	int
}

func NewResourcesClaimManager() *ResourcesClaimManager {
	return &ResourcesClaimManager{
		ResourcesClaims:	map[string]ResourcesClaim{},
	}
}

func (rcm *ResourcesClaimManager) RemoveClaim (claim ResourcesClaim) {
	if rcm.keyExists(claim) {
		delete(rcm.ResourcesClaims, claim.ServerUUID)
		rcm.TotallyUsedVcpus -= claim.Vcpus
		rcm.TotallyUsedDiskGB -= claim.DiskGB
		rcm.TotallyUsedRamMB -= claim.RamMB
	}
}

func (rcm *ResourcesClaimManager) AppendClaim (claim ResourcesClaim) {
	rcm.TotallyUsedVcpus += claim.Vcpus
	rcm.TotallyUsedDiskGB += claim.DiskGB
	rcm.TotallyUsedRamMB += claim.RamMB
	rcm.ResourcesClaims[claim.ServerUUID] = claim
}

func (rcm *ResourcesClaimManager) keyExists (claim ResourcesClaim) bool {
	_, ok := rcm.ResourcesClaims[claim.ServerUUID]
	if ok {
		return true
	}
	return false
}


// TODO: gophercloud gives opportunity to get VMs updated after some time before
// 					servers requests  List ("changes-since")
// TODO: this means we can query periodically about updates and maintain the state
// TODO: of the cluster in our objects, so when the fail signal comes we will have
// TODO: information about all required by VM resources.


// Creates a splice of VMs to evacuate with the rank.
func GetVMsToEvacuate(client *gophercloud.ServiceClient, hostname string) (serversSlice []servers.Server, err error) {
	opts := servers.ListOpts{Host: hostname, AllTenants: true}
	servers.List(client, opts).EachPage(func(page pagination.Page) (bool, error) {
		list, err := servers.ExtractServers(page)
		if err != nil {
			log.Errorf("While extracting servers got: %s", err)
		}
		serversSlice = append(serversSlice, FilterVMsOnEvacuationPolicy(list)...)
		return true, nil
	})

	return serversSlice, nil
}

func FilterVMsOnEvacuationPolicy(serversSlice []servers.Server) (filteredServersSlice []servers.Server){
	for _, server := range serversSlice {
		if server.Metadata["evacuation_policy"] == "Evacuation" {
			filteredServersSlice = append(filteredServersSlice, server)
		}
	}
	return filteredServersSlice
}


type ServerEvacuation struct {
	ServerBefore            servers.Server
	ServerCurrent           servers.Server
	IsEvacuatedSuccessfully bool
	// VM can be scheduled by Healer to different host it actually booted
	// so we need to clear claims of this host
	ScheduledTo             string
}

func NewServerEvacuation (server servers.Server) *ServerEvacuation {
	return &ServerEvacuation{
		ServerBefore: 				server,
		IsEvacuatedSuccessfully:	false,
	}
}

type ByRange []*ServerEvacuation

func (a ByRange) Len() int {return len(a)}

func (a ByRange) Swap(i, j int) {a[i], a[j] = a[j], a[i]}

func (a ByRange) Less(i, j int) bool {
	range_i, ok := a[i].ServerBefore.Metadata["evacuation_range"].(int)
	if !ok {
		range_i = MinEvacuationRangeValue
	}
	range_j, ok := a[j].ServerBefore.Metadata["evacuation_range"].(int)
	if !ok {
		range_j = MinEvacuationRangeValue
	}
	// min should be last - easy pop from slice
	// min value = max priority
	return range_i > range_j
}

// Sorts ascending by evacuation range (0 is the biggest priority)
// If no priority declared - counted as min priority
func SortVMsOnEvacuationRange(evac []*ServerEvacuation) {
	sort.Sort(ByRange(evac))
}


// Run command "nova evacuate .." for a selected VM
// Nova runs command synchronously,so have to wait for
// result.
func (se *ServerEvacuation) Evacuate(client *gophercloud.ServiceClient) error {
	// TODO: 1. create a pool of workers (max size = hypervisors count)
	// TODO: 2. each worker sends Nova command  to evacuate selected VM
	// TODO: 3. worker waits for the result.
	// TODO: 4. worker gets hostname of VM evacuated on (or get it with state?)
	// TODO: 5. worker gets a hypervisor detail for this host and updates Cluster Resourses
	// TODO: STEPS:
	// TODO: 3. create pool of workers
	// TODO: 4. update Cluster Resources
	// TODO: 5. delete claim and evac obj

	return nil
}

func (se *ServerEvacuation) Claim (client *gophercloud.ServiceClient) (ResourcesClaim, error) {
	claim, err := NewResourcesClaim(se.ServerBefore)
	if err != nil {
		log.Error("Could not create a claim for server %s", se.ServerBefore.Name)
	}
	claim.ClaimResources(client)
	return *claim, err
}

func (se *ServerEvacuation) CheckServerEvacuation(client *gophercloud.ServiceClient) error{
	serverNewObj, err := servers.Get(client, se.ServerBefore.ID).Extract()
	if err != nil {
		return err
	}
	if serverNewObj.Status == "ACTIVE" && se.ServerBefore.HostID != serverNewObj.HostID {
		se.IsEvacuatedSuccessfully = true

	}
	se.ServerCurrent = *serverNewObj
	return nil
}