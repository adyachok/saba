package healer

import (
	"sort"

	log "github.com/Sirupsen/logrus"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)
var cluster Cluster

func init() {
	cluster = Cluster{}
}

type Healer struct {
	eventCh              <- chan interface{} // close of this chanel will shutdown healer

	finishedEvacuationCh chan ServerEvacuation
	evacuationCh         chan *ServerEvacuation

	cluster              Cluster
	ResourcesClaims      map[string]ResourcesClaim // key instance UUID

	EvacuationQueue      []ServerEvacuation
	CurrentlyEvacuating  map[string]*ServerEvacuation
}

func NewHealer(event <- chan interface{}) *Healer {
	return &Healer{
		eventCh: event,

		finishedEvacuationCh: make(chan ServerEvacuation),
		evacuationCh: make(chan *ServerEvacuation),

		cluster: Cluster{},
		ResourcesClaims: map[string]ResourcesClaim{},

		EvacuationQueue: []ServerEvacuation{},
		CurrentlyEvacuating: map[string]*ServerEvacuation{}, // key instance UUID
	}
}

func (h *Healer) Shutdown() {
	close(h.evacuationCh)
}

func (h *Healer) Heal(client *gophercloud.ServiceClient) {
	for {
		select {
		// http://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
		case event, ok := <- h.eventCh:
			// TODO:
			if !ok {
				h.Shutdown()
				break
			}
			// TODO: fix this event check
			if event == "fail" {
				// TODO: update evacuation queue:
				// TODO: 1. get VMs to evacuate
				// TODO: 2. range VMs
				// We have evacuation order from max to min - but min value has
				// biggest evacuation priority, so we have to reverse
				var evacuationOrderReversed []ServerEvacuation
				copy(evacuationOrderReversed, h.CurrentlyEvacuating)
				sort.Sort(sort.Reverse(evacuationOrderReversed))

				for _, server := range evacuationOrderReversed {
					err := cluster.UpdateAvailableClusterResources(client)
					if err != nil {
						log.Errorf("Error updating cluster stare: %s", err)
						continue
					}
					claim, err  := h.claimResourcesForServer(client, server.ServerBefore)
					if err != nil {
						log.Errorf("Error claiming resources: %s", err)
						continue
					}
					// TODO: update cluster state
					// TODO: count available resources = cluster res - claimed
					// TODO: if cluster has enogh resources
					select {
					case h.evacuationCh <- server:
						// We pop send to evacuation server from wait evacuation queue
						h.EvacuationQueue = h.EvacuationQueue[:len(h.EvacuationQueue) - 1]
						h.CurrentlyEvacuating[server.ServerBefore.ID] = &server

					}
				}
			}
		case server := <- h.finishedEvacuationCh:
			// TODO: check success/error, log
			delete(h.CurrentlyEvacuating, server.ServerBefore.ID)
		}

	}
}

func (h *Healer) claimResourcesForServer(client *gophercloud.ServiceClient, server servers.Server) (claim *ResourcesClaim, err error){
	claim, err = NewResourcesClaim(server)
	if err != nil {
		return nil, err
	}
	err = claim.ClaimResources(client)
	if err != nil {
		return nil, err
	}
	return claim, err
}