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
are stored in Evac_Q of the Healer object.
 
When filtering and sorting finished, process executes scheduling logic
on every VM. Scheduling applies simple filters on available CPU, RAM and
Hard Drive memory. If VM fits - it is routed to Scheduled_Q, if not - to
the Failed_Q. 