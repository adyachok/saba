package healer

import (
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/adyachok/bacsi/openstack/v2/services"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestPoolEvacuate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)
	qm := NewQueueManager()

	client := client.ServiceClient()
	healer.prepareVMsEvacuation(client, qm)

	th.AssertEquals(t, 2, len(qm.Scheduled_Q))


	resultsCh := make(chan *EvacContainer)

	dispatcher := NewDispatcher(client, resultsCh)
	dispatcher.activate(qm)

	ecSlice := []*EvacContainer{}

	for i:=1; i < 3; i++ {
		res := <- resultsCh
		ecSlice = append(ecSlice, res)
		th.AssertEquals(t, "accepted", res.State)
	}

	th.AssertEquals(t, 2, len(ecSlice))
	dispatcher.shutdown()
}

func TestPoolCheckEvacuation(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)
	qm := NewQueueManager()

	client := client.ServiceClient()
	healer.prepareVMsEvacuation(client, qm)

	th.AssertEquals(t, 2, len(qm.Scheduled_Q))


	resultsCh := make(chan *EvacContainer)
	t.Log(qm)
	dispatcher := NewDispatcher(client, resultsCh)
	go dispatcher.activate(qm)


	for i:=1; i < 3; i++ {
		res := <- resultsCh
		th.AssertEquals(t, "accepted", res.State)
		healer.processAccepedContainer(res, qm)
		t.Log(qm)
	}
	//th.AssertEquals(t, 2, l)
	dispatcher.shutdown()
}