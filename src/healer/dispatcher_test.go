package healer

import (
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/adyachok/bacsi/openstack/v2/services"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestPool_Run(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)

	client := client.ServiceClient()
	healer.prepareVMsEvacuation(client)

	th.AssertEquals(t, 2, len(healer.queryManager.Scheduled_Q))


	resultsCh := make(chan *EvacContainer)

	dispatcher := NewDispatcher(client, resultsCh)
	dispatcher.activate(healer.queryManager.Scheduled_Q, healer.queryManager.Accepted_Q)

	ecSlice := []*EvacContainer{}

	for i:=1; i < 3; i++ {
		res := <- resultsCh
		ecSlice = append(ecSlice, res)
		th.AssertEquals(t, "accepted", res.State)
	}

	th.AssertEquals(t, 2, len(ecSlice))
	dispatcher.shutdown()
}
