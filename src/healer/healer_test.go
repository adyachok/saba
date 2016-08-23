package healer

import (
	"time"
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	"github.com/adyachok/bacsi/openstack/v2/services"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestHeal(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	services.HandleServiceForceDownSuccessfully(t)
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)
	HandleServersSuccessfully(t)

	evCh := make(chan interface{})
	healer := NewHealer(evCh)

	go func (){
		event := "join"
		evCh <- event
		time.Sleep(1 * time.Second)
		event = "fail"
		evCh <- event
		time.Sleep(1 * time.Second)
		close(evCh)
	}()

	client := client.ServiceClient()
	healer.Heal(client)
}