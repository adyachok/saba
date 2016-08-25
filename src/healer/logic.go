package healer

import (
	"errors"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
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
	c.Resources = make(map[string]HypervisorFreeResources)

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
