package saba

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adyachok/bacsi/openstack/v2/hypervisors"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)


func TestGetAvailableClusterResources(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	hypervisors.HandleHypervisorsDetailsSuccessfully(t)

	clusterExpected := Cluster{
		Resources: map[string]HypervisorFreeResources{
			"compute-0-3.domain.tld": HypervisorFreeResources{"compute-0-3.domain.tld", 2, "up", "enabled", 24, "192.168.2.28", 414, 633, 5120,},
			"compute-0-1.domain.tld": HypervisorFreeResources{"compute-0-1.domain.tld", 8, "up", "enabled", 22, "192.168.2.25", 364, 601, 2048,},
			"compute-0-2.domain.tld": HypervisorFreeResources{"compute-0-2.domain.tld", 14, "up", "enabled", 30, "192.168.2.23", 524, 521, 44032,},
		},
	}

	cluster := Cluster{}
	err := GetAvailableClusterResources(client.ServiceClient(), &cluster)

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, clusterExpected, cluster)
}

func TestSortVMsOnEvacuationRange(t *testing.T) {
	testServersList := make([]servers.Server, len(ServersList))
	copy(testServersList, ServersList)
	SortVMsOnEvacuationRange(testServersList)

	if len(ServersList) != 3 {
		t.Errorf("Expected 3 servers, saw %d", len(testServersList))
	}
	th.AssertDeepEquals(t, ServersListSortedByRangeExpected, testServersList)
}

func TestFilterVMsOnEvacuationPolicy(t *testing.T) {
	testServersList := make([]servers.Server, len(ServersList))
	copy(testServersList, ServersList)
	filtered := FilterVMsOnEvacuationPolicy(testServersList)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 servers, saw %d", len(filtered))
	}
}

func TestGetVMsToEvacuate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ServersListBody)
	})

	serversSlice, err := GetVMsToEvacuate(client.ServiceClient(), "compute-0-4.domain.tld")

	if len(serversSlice) != 2 {
		t.Errorf("Expected 2 servers, saw %d", len(serversSlice))
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ServersListFilteredByEvacuationPolicyExpected, serversSlice)
}