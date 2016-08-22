**Healer**
This object should contain next fields:
    1. Evaq_Q = []VM
    2. FailedEvaq_Q = []VM
    3. Scheduled_M = map[vm_id] compute_host_id
    4. Claims_M = map[compute_host_id] Claim

**States**
Watch file: ~/Pictures/Proj/DSC_0342.JPG

On future: add "join" event processing, to trigger evacuation of not 
evacuated VMs.