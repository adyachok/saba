package saba

import (
	"sort"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/pagination"
	log "github.com/Sirupsen/logrus"
)


const (
	MaxUint = ^uint(0)
 	MinEvacuationRangeValue = int(MaxUint >> 1) // MaxInt will be min priority
)


type HypervisorFreeResources struct {
	HypervisorHostname string		`mapstructure:"hypervisor_hostname"`
	Id int
	State string
	Status string
	FreeVcpus int16					`mapstructure:"free_vcpus"` 	// Vcpus - VcpuUsed
	HostIP string 					`mapstructure:"host_ip"`
	FreeDiskGB int16				`mapstructure:"free_disk_gb"`
	// disk_available_least = free_disk_gb - disk_overcommit_size
	// disk_overcommit_size = virtual size of disks of all instance instance - used disk size of all instances
	DistAvailableLeast int16		`mapstructure:"disk_available_least"`
	FreeRamMB int32					`mapstructure:"free_ram_mb"`
}

type Cluster struct {
	Resources map[string]HypervisorFreeResources
}

// Creates a slice of every host resources
func GetAvailableClusterResources(client *gophercloud.ServiceClient, cluster *Cluster) error {
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
		cluster.Resources = make(map[string]HypervisorFreeResources)
		cluster.Resources[hypervisorDetails.HypervisorHostname] = res
	}
	return nil
}

// Creates claim for resources booting VM
func claimResources(){

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

type ByRange []servers.Server

func (a ByRange) Len() int {return len(a)}

func (a ByRange) Swap(i, j int) {a[i], a[j] = a[j], a[i]}

func (a ByRange) Less(i, j int) bool {
	range_i, ok := a[i].Metadata["evacuation_range"].(int)
	if !ok {
		range_i = MinEvacuationRangeValue
	}
	range_j, ok := a[j].Metadata["evacuation_range"].(int)
	if !ok {
		range_j = MinEvacuationRangeValue
	}
	return range_i < range_j
}

// Sorts ascending by evacuation range (0 is the biggest priority)
// If no priority declared - counted as min priority
func SortVMsOnEvacuationRange(serversSlice []servers.Server) {
	sort.Sort(ByRange(serversSlice))
}

// Run command "nova evacuate .." for a selected VM
// Nova runs command synchronously,so have to wait for
// result.
func Evacuate() error {
	return nil
}