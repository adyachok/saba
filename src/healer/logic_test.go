package healer

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
	err := cluster.UpdateAvailableClusterResources(client.ServiceClient())

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, clusterExpected, cluster)
}

func TestSortVMsOnEvacuationRange(t *testing.T) {
	var testServersList []*EvacContainer
	for _, serv := range ServersList {
		testServersList = append(testServersList, NewEvacContainer(serv))
	}

	SortVMsOnEvacuationRange(testServersList)

	if len(testServersList) != 3 {
		t.Errorf("Expected 3 servers to Evacuate, saw %d", len(testServersList))
	}

	th.AssertDeepEquals(t, ServersListSortedByRangeExpected, testServersList)
	th.AssertEquals(t, 0, ServersListSortedByRangeExpected[2].ServerBefore.Metadata["evacuation_range"])
	th.AssertEquals(t, 100, ServersListSortedByRangeExpected[1].ServerBefore.Metadata["evacuation_range"])
	th.AssertEquals(t, nil, ServersListSortedByRangeExpected[0].ServerBefore.Metadata["evacuation_range"])
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
	HandleServersSuccessfully(t)

	serversSlice, err := GetVMsToEvacuate(client.ServiceClient(), "compute-0-4.domain.tld")

	if len(serversSlice) != 2 {
		t.Errorf("Expected 2 servers, saw %d", len(serversSlice))
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ServersListFilteredByEvacuationPolicyExpected, serversSlice)
}

func TestCheckServerEvacuationSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ServerDerpSuccessfulEvacuationRespBody)
	})
	se := NewEvacContainer(ServerDerp)
	err := se.CheckServerEvacuation(client.ServiceClient())

	if !se.IsEvacuatedSuccessfully {
		t.Errorf("Expected server Derp to be evacuated")
	}

	th.AssertNoErr(t, err)
}

func TestCheckServerEvacuationFail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, ServerDerpFailedEvacuationRespBody)
	})
	se := NewEvacContainer(ServerDerp)

	err := se.CheckServerEvacuation(client.ServiceClient())

	if se.IsEvacuatedSuccessfully {
		t.Errorf("Expected server Derp not to be evacuated")
	}

	th.AssertNoErr(t, err)
}

func TestClaimResources(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleFlavorGetSuccessfully(t)

	claim, err := NewResourcesClaim(ServerDerp)
	if err != nil {
		t.Errorf("Resources claim object creation error: %s", err)
	}
	err = claim.ClaimResources(client.ServiceClient())
	if err != nil {
		t.Errorf("Resources claim error: %s", err)
	}
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ResourceClaimExpected, claim)
}

func TestClaimsCount(t *testing.T) {
	rcm := NewResourcesClaimManager()
	rcm.AppendClaim(*ResourceClaimExpected)

	th.AssertEquals(t, 1, len(rcm.ResourcesClaims))
	th.AssertEquals(t, ResourceClaimExpected.DiskGB, rcm.TotallyUsedDiskGB)
	th.AssertEquals(t, ResourceClaimExpected.RamMB, rcm.TotallyUsedRamMB)
	th.AssertEquals(t, ResourceClaimExpected.Vcpus, rcm.TotallyUsedVcpus)

	rcm.RemoveClaim(*ResourceClaimExpected)
	th.AssertEquals(t, 0, len(rcm.ResourcesClaims))
	th.AssertEquals(t, 0, rcm.TotallyUsedDiskGB)
	th.AssertEquals(t, 0, rcm.TotallyUsedRamMB)
	th.AssertEquals(t, 0, rcm.TotallyUsedVcpus)

}