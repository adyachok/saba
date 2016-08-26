**Healer**
This object should contain next fields:
    1. Evaq_Q = []VM
    2. FailedEvaq_Q = []VM
    3. Scheduled_M = map[vm_id] compute_host_id
    4. Claims_M = map[compute_host_id] Claim
    
    
Unite whole logic in two channels:
    1. taskCh
    2. resultCh
    
Logic:
        case event := <-eventCh        
        case taskCh <-  Assert_Q (prio 1) or Evac_Q (prio 2)
        case server := <- resultCh

Also add states to EvacContainer:
    state : on from ["scheduled", "accepted", "finished", "failed"]

**States**
Watch file: ~/Pictures/Proj/DSC_0342.JPG

On future: add "join" event processing, to trigger evacuation of not 
evacuated VMs.