package healer

import (
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/adyachok/bacsi/openstack/v2/services"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestPoolFailEvacuate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)
	ServerDerpFailedEvacuationHandler(t)
	ServerHerpFailedEvacuationHandler(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)
	qm := NewQueueManager()

	client := client.ServiceClient()
	healer.prepareVMsEvacuation(client, qm)

	th.AssertEquals(t, 2, len(qm.Scheduled_Q))


	resultsCh := make(chan *EvacContainer)

	dispatcher := NewDispatcher(client, resultsCh)
	go dispatcher.activate(qm)

	for i:=4; i > 0; i-- {
		res := <- resultsCh
		switch {
		case res.State == "accepted":
			healer.processAcceptedContainer(res, qm)
		case res.State == "failed":
			healer.processFailedContainer(*res, qm)
		}
	}
	th.AssertEquals(t, 2, len(healer.FailedEvac_Q))

	dispatcher.shutdown()
}


func TestPoolSuccessfulEvacuate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)
	ServerDerpSuccessfulEvacuationHandler(t)
	ServerHerpSuccessfulEvacuationHandler(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)
	qm := NewQueueManager()

	client := client.ServiceClient()
	healer.prepareVMsEvacuation(client, qm)

	th.AssertEquals(t, 2, len(qm.Scheduled_Q))


	resultsCh := make(chan *EvacContainer)

	dispatcher := NewDispatcher(client, resultsCh)
	go dispatcher.activate(qm)

	var finishedCounter int

	for i:=4; i > 0; i-- {
		res := <- resultsCh
		switch {
		case res.State == "accepted":
			healer.processAcceptedContainer(res, qm)
		case res.State == "finished":
			healer.processFinishedContainer(*res, qm)
			finishedCounter++
		}
	}

	th.AssertEquals(t, 2, finishedCounter)

	dispatcher.shutdown()
}