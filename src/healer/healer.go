package healer

import (
	"fmt"

	"github.com/adyachok/bacsi/openstack/v2/services"
	log "github.com/Sirupsen/logrus"
	"github.com/rackspace/gophercloud"
)

const (
	MaxRetries = 4
)

var cluster Cluster

func init() {
	cluster = Cluster{}
}

type Healer struct {
	// Close of eventCh will shutdown healer
	eventCh              	<- chan interface{}
	finishedEvacuationCh 	chan ServerEvacuation
	evacuationCh         	chan *ServerEvacuation
	cluster              	Cluster
	// mapping between hypervisor hostname {key} and claims to this hypervisor
	Claims_M             	map[string]*ResourcesClaimManager
	Evac_Q               	[] *ServerEvacuation
	Scheduled_Q          	[] *ServerEvacuation
	FailedEvac_Q         	[] ServerEvacuation
}

func NewHealer(event <- chan interface{}) *Healer {
	return &Healer{
		eventCh: 				event,
		finishedEvacuationCh: 	make(chan ServerEvacuation),
		evacuationCh: 			make(chan *ServerEvacuation),
		cluster: 				Cluster{},
		// Mapping of compute id {key} and slices of claimed resources
		Claims_M: 				map[string] *ResourcesClaimManager{},
		Evac_Q: 				[] *ServerEvacuation{},
		Scheduled_Q: 			[] *ServerEvacuation{},
	}
}

func (h *Healer) Shutdown() {
	close(h.evacuationCh)
}

func (h *Healer) Heal(client *gophercloud.ServiceClient) {
	for {
		select {
		// http://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
			case event, ok := <-h.eventCh:
				if !ok {
					h.Shutdown()
					break
				}
				evt, ok := event.(string)
				if !ok {
					continue
				}
				switch {
					case evt == "fail":
						hostname := h.getFailedHostname()
						err := h.forceDownServiceWithRetry(client, hostname)
						// We append new VMs only if compute service was successfully
						// forced down, in other case VMs won't be successfully evacuated.
						if err == nil {
							h.Evac_Q = append(h.Evac_Q, h.updateEvacuationQueueWithRetry(client, hostname)...)
							// We have evacuation order from max to min - but min value has
							// biggest evacuation priority, so we have to reverse
							SortVMsOnEvacuationRange(h.Evac_Q)
						}

						err = cluster.UpdateAvailableClusterResources(client)
						if err != nil {
							log.Errorf("Error updating cluster available resources: %s", err)
						}

						for i := len(h.Evac_Q)-1; i >= 0; i-- {
							server :=  h.Evac_Q[i]
							h.Evac_Q = append(h.Evac_Q[:i], h.Evac_Q[i+1:]...)


							claim := h.claimResourcesWithRetry(client, server)
							if claim == nil {
								h.FailedEvac_Q = append(h.FailedEvac_Q, *server)
								continue
							}
							err = h.schedule(server, *claim)
							if err != nil {
								// No available resources found for VM... skipping
								continue
							}

							select {
							case h.evacuationCh <- server:
								h.Scheduled_Q = append(h.Scheduled_Q, server)
							}
						}
					case evt == "join":
						// TODO:
						log.Info("Server joined")
					default:
						log.Infof("Got unexpected event %s", event)

				}

			case server := <-h.finishedEvacuationCh:
				for idx, server_ := range h.Scheduled_Q {
					if server_.ServerBefore.ID == server.ServerBefore.ID {
						// Remove from scheduled  queue
						h.Scheduled_Q = append(h.Scheduled_Q[:idx], h.Scheduled_Q[idx+1:]...)
						break
					}
				}
				if !server.IsEvacuatedSuccessfully {
					// If server wasn't evacuated successfully we add it to the
					// failed queue
					log.Errorf("Server %s wasn't evacuated successfully", server.ServerBefore.ID)
					h.FailedEvac_Q = append(h.FailedEvac_Q, server)
				}
			log.Infof("Server %s was evacuated successfully", server.ServerBefore.ID)
		}
	}
}

func (h *Healer) getFailedHostname() string {
	// TODO:
	return "compute-0-1"
}

func (h *Healer) updateEvacuationQueueWithRetry(client *gophercloud.ServiceClient, hostname string) []*ServerEvacuation {
	for i := 0; i <= MaxRetries; i++ {
		vms, err := GetVMsToEvacuate(client, hostname)
		if err != nil {
			log.Errorf("Got error: %s during get VMs to evacuate", err)
		}else {
			vms_ := []*ServerEvacuation{}
			for _, vm := range vms {
				se := NewServerEvacuation(vm)
				vms_ = append(vms_, se)
			}
			return vms_
		}
	}
	return nil
}

// Find out appropriate host to evacuate an instance or add instance to the
// failing evacuation list
func (h *Healer) schedule(instance *ServerEvacuation, claim ResourcesClaim) error {
	for hostname, resources := range h.cluster.Resources {
		if h.filterResources(claim, hostname, resources) {
			h.Claims_M[hostname].AppendClaim(claim)
			instance.ScheduledTo = hostname
			return nil
		}
	}
	err := fmt.Errorf("No valid host found for instance %s", claim.ServerUUID)
	log.Error(err)
	h.FailedEvac_Q = append(h.FailedEvac_Q, *instance)
	return err
}

// Simple dumb filtering.
func (h *Healer) filterResources(claim ResourcesClaim, hostname string, resources HypervisorFreeResources) bool {
	claimsManager := h.Claims_M[hostname]

	if resources.FreeVcpus - int16(claimsManager.TotallyUsedVcpus) <= int16(claim.Vcpus) {
		return false
	}
	if resources.FreeDiskGB - int16(claimsManager.TotallyUsedDiskGB) <= int16(claim.DiskGB) {
		return false
	}
	if resources.FreeRamMB - int32(claimsManager.TotallyUsedRamMB) <= int32(claim.RamMB) {
		return false
	}
	return true
}

func (h *Healer) forceDownServiceWithRetry(client *gophercloud.ServiceClient, hostname string) (exec_err error) {
	for i:=0; i <= MaxRetries; i++ {
		result, err := services.ForceDown(client, "nova-compute", hostname).Extract()
		if err != nil {
			log.Errorf("Got error: %s during nova-compute force down attempt on host %s", err, hostname)
		}
		if !result.ForcedDown {
			err = fmt.Errorf("Cannot force down nova-compute service on host %s", hostname)
			log.Error(err)
		} else {
			return nil
		}
		exec_err = err
	}
	return exec_err
}

func (h *Healer) claimResourcesWithRetry(client *gophercloud.ServiceClient, server *ServerEvacuation) *ResourcesClaim {
	for i:=0; i <= MaxRetries; i++ {
		claim, err := server.Claim(client)
		if err != nil {
			log.Errorf("Got error: %s during claiming resources for instance %s", err, server.ServerBefore.ID)
		}else {
			return &claim
		}
	}
	return nil
}