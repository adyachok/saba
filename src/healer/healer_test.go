package healer

import (
	"time"
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/adyachok/bacsi/openstack/v2/services"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestFilterResources(t *testing.T) {
	evCh := make(chan interface{})
	healer := NewHealer(evCh)

	instance := NewServerEvacuation(ServerDerp)
	result := healer.filterResources(*ResourceClaimExpected, instance.ServerBefore.HostID, HypervisorFreeResources_1)
	th.AssertEquals(t, true, result)
}


func TestHealerSchedule(t *testing.T) {
	evCh := make(chan interface{})
	healer := NewHealer(evCh)

	healer.cluster.Resources = map[string]HypervisorFreeResources{}
	healer.cluster.Resources["compute-0-1"] = HypervisorFreeResources_1

	instance := NewServerEvacuation(ServerDerp)
	healer.schedule(instance, *ResourceClaimExpected)

	th.AssertEquals(t, 1, len(healer.Claims_M))
	th.AssertEquals(t, "compute-0-1", instance.ScheduledTo)
}

func TestHeal(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)

	VMsToEvacuate := []*ServerEvacuation{}

	go func (){
		event := "join"
		evCh <- event

		time.Sleep(1 * time.Second)
		event = "fail"

		evCh <- event
		time.Sleep(1 * time.Second)

		vm := <-healer.evacuationCh
		VMsToEvacuate = append(VMsToEvacuate, vm)

		vm = <-healer.evacuationCh
		VMsToEvacuate = append(VMsToEvacuate, vm)

		time.Sleep(1 * time.Second)
		close(evCh)
	}()

	client := client.ServiceClient()
	healer.Heal(client)

	th.AssertEquals(t, 2, len(VMsToEvacuate))
}

