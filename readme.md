Healer
======
The project takes care about your cloud. It traks for the states of 
computes in the cluster and in case of failing on of them starts to 
evacuate virtual machines to another, healthy part of the cluster.


Interaction overview
====================

Healer - Dispatcher
-------------------
![Alt text](pict/DSC_0371.png?raw=true "Healer - Dispatcher interaction overview")

Healer has to react on host failed or joined event.

On compute failed event process tries to get all VMs hosted on failed
compute. After getting VMs - they are filtered accordingly to evacuation
policy (specified in VM metadata), as well as, sorted accordingly to
evacuation priority (specified in VM metadata). Filtered and sorted VMs
are stored in **Evac_Q** of the Healer object.
 
When filtering and sorting finished, process executes scheduling logic
on every VM. Scheduling applies simple filters on available CPU, RAM and
Hard Drive memory. If VM fits - it is routed to **Scheduled_Q**, if not 
- to the **Failed_Q**.
 
Scheduled_Q belongs to QueueManager object which is shared between Healer
and Dispatcher. Dispatcher, in active state, periodically checks for
objects in Scheduled_Q and Accepted_Q, and if one of the queues is not
empty, it pops and sends an object to task execution pool.

During task execution accordingly to the inner logic of EvacContainer, it
receives predifined state.Task execution pool executes a task and sends
object to the Healed object.

Healer object accordingly to inner logic checks state of an object and
process it.

On last EvacContainer comes evacuation queue can be updated accordingly
to resources available.


EvacContainer
-------------

EvacContainer contains information about specified VM. 
Field **ServerBefore** contains Server object with information before
evacuation.
Field **ServerCurrent** contains Server object with information after
server evacuation.
Fiels **State** contains information about EvacContainer state accordingly
to evacuation process. There are next states: 
    1. "scheduled" 
    2. "accepted" 
    3. "finished" 
    4. "failed"

