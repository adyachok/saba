package healer

import (
	"sort"
	"github.com/rackspace/gophercloud"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	log "github.com/Sirupsen/logrus"
	"github.com/rackspace/gophercloud/pagination"
)

type EvacContainer struct {
	Id            string
	ServerBefore  servers.Server
	ServerCurrent servers.Server
	// "scheduled", "accepted", "finished", "failed"
	State         string
	// VM can be scheduled by Healer to different host it actually booted
	// so we need to clear claims of this host
	ScheduledTo   string
}

func (h *EvacContainer) Task(client *gophercloud.ServiceClient) {
	switch {
	case h.State == "scheduled":
		h.Evacuate(client)
	case h.State == "accepted":
		h.CheckServerEvacuation(client)
	}
}

func NewEvacContainer(server servers.Server) *EvacContainer {
	return &EvacContainer{
		Id:							server.ID,
		ServerBefore: 				server,
	}
}

type ByRange []*EvacContainer

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
func SortVMsOnEvacuationRange(evac []*EvacContainer) {
	sort.Sort(ByRange(evac))
}


// Run command "nova evacuate .." for a selected VM
// Nova runs command synchronously,so have to wait for
// result.
func (se *EvacContainer) Evacuate(client *gophercloud.ServiceClient) error {
	// TODO: 1. create a pool of workers (max size = hypervisors count)
	// TODO: 2. each worker sends Nova command  to evacuate selected VM
	// TODO: 3. worker waits for the result.
	// TODO: 4. worker gets hostname of VM evacuated on (or get it with state?)
	// TODO: 5. worker gets a hypervisor detail for this host and updates Cluster Resourses
	// TODO: STEPS:
	// TODO: 3. create pool of workers
	// TODO: 4. update Cluster Resources
	// TODO: 5. delete claim and evac obj

	se.State = "accepted" //If accepted or fail in case of error
	return nil
}

func (se *EvacContainer) Claim (client *gophercloud.ServiceClient) (ResourcesClaim, error) {
	claim, err := NewResourcesClaim(se.ServerBefore)
	if err != nil {
		log.Error("Could not create a claim for server %s", se.ServerBefore.Name)
	}
	claim.ClaimResources(client)
	return *claim, err
}

func (se *EvacContainer) CheckServerEvacuation(client *gophercloud.ServiceClient) error{
	serverNewObj, err := servers.Get(client, se.Id).Extract()
	if err != nil {
		return err
	}
	if serverNewObj.Status == "ACTIVE" && se.ServerBefore.HostID != serverNewObj.HostID {
		se.State = "finished"

	}
	if serverNewObj.Status == "ERROR" {
		se.State = "failed"

	}
	se.ServerCurrent = *serverNewObj
	return nil
}

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
