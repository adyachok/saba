package healer

import (
	"fmt"
	"sync"

	"github.com/adyachok/bacsi/openstack/v2/services"
	log "github.com/Sirupsen/logrus"
	"github.com/rackspace/gophercloud"

)

const (
	MaxRetries = 4
)

type QueueManager struct {
	lock		    sync.RWMutex
	Scheduled_Q    [] *EvacContainer
	Accepted_Q	   [] *EvacContainer
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Scheduled_Q:    [] *EvacContainer{},
		Accepted_Q:	   [] *EvacContainer{},
	}
}

type Healer struct {
	// Close of eventCh will shutdown healer
	eventCh      <- chan interface{}
	resultCh     chan *EvacContainer
	cluster      Cluster
	// mapping between hypervisor hostname {key} and claims to this hypervisor
	Claims_M     map[string]*ResourcesClaimManager
	Evac_Q       [] *EvacContainer
	FailedEvac_Q [] EvacContainer
}

func NewHealer(event <- chan interface{}) *Healer {

	return &Healer{
		eventCh: 		event,
		resultCh: 		make(chan *EvacContainer),
		cluster: 		Cluster{},
		// Mapping of compute id {key} and slices of claimed resources
		Claims_M: 		map[string] *ResourcesClaimManager{},
		Evac_Q: 		[] *EvacContainer{},
		FailedEvac_Q:   [] EvacContainer{},
	}
}

func (h *Healer) Shutdown(dispatcher *Dispatcher) {
	dispatcher.shutdown()
}

func (h *Healer) Heal(client *gophercloud.ServiceClient) {
	dispatcher := NewDispatcher(client, h.resultCh)
	queueManager := NewQueueManager()
	go dispatcher.activate(queueManager)

	for {
		select {
		// http://stackoverflow.com/questions/13666253/breaking-out-of-a-select-statement-when-all-channels-are-closed
			case event, ok := <-h.eventCh:
				if !ok {
					h.Shutdown(dispatcher)
					// TODO: wait for subprocesses to stop
					return
				}
				evt, ok := event.(string)
				if !ok {
					continue
				}
				switch {
					case evt == "fail":
						h.prepareVMsEvacuation(client, queueManager)

					case evt == "join":
						// TODO:
						log.Info("Server joined")
					default:
						log.Infof("Got unexpected event %s", event)

				}

			case container := <-h.resultCh:
				switch {
					case container.State == "accepted":
						h.processAcceptedContainer(container, queueManager)
					case container.State == "finised":
						h.processFinishedContainer(*container, queueManager)
						h.Reschedule(client, queueManager)
					case container.State == "failed":
						h.processFailedContainer(*container, queueManager)
						h.Reschedule(client, queueManager)
				}
		}
	}
}

func (h *Healer) prepareVMsEvacuation(client *gophercloud.ServiceClient, qm *QueueManager) {
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

	h.Schedule(client, qm)
}


func (h *Healer) getFailedHostname() string {
	// TODO:
	return "compute-0-1"
}

func (h *Healer) updateEvacuationQueueWithRetry(client *gophercloud.ServiceClient, hostname string) []*EvacContainer {
	for i := 0; i <= MaxRetries; i++ {
		vms, err := GetVMsToEvacuate(client, hostname)
		if err != nil {
			log.Errorf("Got error: %s during get VMs to evacuate", err)
		}else {
			vms_ := []*EvacContainer{}
			for _, vm := range vms {
				se := NewEvacContainer(vm)
				vms_ = append(vms_, se)
			}
			return vms_
		}
	}
	return nil
}

func (h *Healer) Schedule(client *gophercloud.ServiceClient, qm *QueueManager) {
	err := h.cluster.UpdateAvailableClusterResources(client)
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
		err = h.ScheduleVM(server, *claim)
		if err != nil {
			// No available resources found for VM... skipping
			continue
		}

		qm.lock.RLock()
		qm.Scheduled_Q = append(qm.Scheduled_Q, server)
		qm.lock.RUnlock()
	}
}

func (h *Healer) Reschedule(client *gophercloud.ServiceClient, qm *QueueManager) {
	if len(qm.Scheduled_Q) == 0 {
		h.Schedule(client, qm)
	}
}

// Find out appropriate host to evacuate an instance or add instance to the
// failing evacuation list
func (h *Healer) ScheduleVM(container *EvacContainer, claim ResourcesClaim) error {
	for hostname, resources := range h.cluster.Resources {
		if h.filterResources(claim, hostname, resources) {
			h.Claims_M[hostname].AppendClaim(claim)
			container.ScheduledTo = hostname
			container.State = "scheduled"
			return nil
		}
	}
	err := fmt.Errorf("No valid host found for instance %s", claim.ServerUUID)
	log.Error(err)
	container.State = "failed"
	h.FailedEvac_Q = append(h.FailedEvac_Q, *container)
	return err
}

// Simple dumb filtering.
func (h *Healer) filterResources(claim ResourcesClaim, hostname string, resources HypervisorFreeResources) bool {
	claimsManager, ok := h.Claims_M[hostname]
	if !ok {
		claimsManager = NewResourcesClaimManager()
		h.Claims_M[hostname] = claimsManager
	}
	free_cpus := resources.FreeVcpus - int16(claimsManager.TotallyUsedVcpus)
	if free_cpus <= int16(claim.Vcpus) {
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

func (h *Healer) claimResourcesWithRetry(client *gophercloud.ServiceClient, server *EvacContainer) *ResourcesClaim {
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

func (h *Healer) processAcceptedContainer(container *EvacContainer, qm *QueueManager) {
	qm.lock.RLock()
	qm.Accepted_Q = append(qm.Accepted_Q, container)
	qm.lock.RUnlock()
}


func (h *Healer) processFinishedContainer(container EvacContainer, qm *QueueManager) {
	log.Infof("Server %s was evacuated successfully", container.ServerBefore.ID)
}

func (h *Healer) processFailedContainer(container EvacContainer, qm *QueueManager) {
	// If server wasn't evacuated successfully we add it to the
	// failed queue
	log.Errorf("Server %s wasn't evacuated successfully", container.ServerBefore.ID)
	h.FailedEvac_Q = append(h.FailedEvac_Q, container)
}