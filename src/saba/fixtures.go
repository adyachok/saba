package saba

import "github.com/rackspace/gophercloud/openstack/compute/v2/servers"


const ServersListBody = `
{
	"servers": [
		{
			"OS-DCF:diskConfig": "MANUAL",
			"OS-EXT-STS:power_state": 1,
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"accessIPv4": "",
			"accessIPv6": "",
			"addresses": {},
			"created": "2014-09-25T13:04:41Z",
			"flavor": {
				"id": "performance1-8",
				"links": [
					{
						"href": "https://dfw.servers.api.rackspacecloud.com/111111/flavors/performance1-8",
						"rel": "bookmark"
					}
				]
			},
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"id": "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"image": {
				"id": "255df5fb-e3d4-45a3-9a07-c976debf7c14",
				"links": [
					{
						"href": "https://dfw.servers.api.rackspacecloud.com/111111/images/255df5fb-e3d4-45a3-9a07-c976debf7c14",
						"rel": "bookmark"
					}
				]
			},
			"key_name": "",
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/111111/servers/59818cee-bc8c-44eb-8073-673ee65105f7",
					"rel": "self"
				},
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/111111/servers/59818cee-bc8c-44eb-8073-673ee65105f7",
					"rel": "bookmark"
				}
			],
			"metadata": {
			"evacuation_policy": "Evacuation",
				"evacuation_range": 100
			},
			"name": "derp",
			"progress": 0,
			"status": "ACTIVE",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"updated": "2014-09-25T13:04:49Z",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd"
		},
		{
			"OS-DCF:diskConfig": "MANUAL",
			"OS-EXT-STS:power_state": 1,
			"OS-EXT-STS:task_state": null,
			"OS-EXT-STS:vm_state": "active",
			"accessIPv4": "",
			"accessIPv6": "",
			"addresses": {},
			"created": "2014-09-25T13:10:02Z",
			"flavor": {
				"id": "performance1-2",
				"links": [
					{
						"href": "https://dfw.servers.api.rackspacecloud.com/111111/flavors/performance1-2",
						"rel": "bookmark"
					}
				]
			},
			"hostId": "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
			"id": "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
						"security_groups": [
				{
					"name": "default"
				}
			],
			"image": {
				"id": "bb02b1a3-bc77-4d17-ab5b-421d89850fca",
				"links": [
					{
						"href": "https://dfw.servers.api.rackspacecloud.com/111111/images/bb02b1a3-bc77-4d17-ab5b-421d89850fca",
						"rel": "bookmark"
					}
				]
			},
			"key_name": "",
			"links": [
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/v2/111111/servers/25f1c7f5-e00a-4715-b355-16e24b2f4630",
					"rel": "self"
				},
				{
					"href": "https://dfw.servers.api.rackspacecloud.com/111111/servers/25f1c7f5-e00a-4715-b355-16e24b2f4630",
					"rel": "bookmark"
				}
			],
			"metadata": {
				"evacuation_policy": "Evacuation",
				"evacuation_range": 0
			},
			"name": "herp",
			"progress": 0,
			"status": "ACTIVE",
			"tenant_id": "fcad67a6189847c4aecfa3c81a05783b",
			"updated": "2014-09-25T13:10:10Z",
			"user_id": "9349aff8be7545ac9d2f1d00999a23cd"
		},
		{
			"OS-EXT-STS:task_state": null,
			"addresses": {
				"private": [
					{
						"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:80:35:77",
						"version": 4,
						"addr": "10.33.160.11",
						"OS-EXT-IPS:type": "fixed"
					}
				]
			},
			"links": [
				{
					"href": "http://192.168.2.21:8774/v2/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
					"rel": "self"
				},
				{
					"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
					"rel": "bookmark"
				}
			],
			"image": {
				"id": "d74d06f2-61f5-4d02-b568-9d5002cf27f4",
				"links": [
					{
						"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/images/d74d06f2-61f5-4d02-b568-9d5002cf27f4",
						"rel": "bookmark"
					}
				]
			},
			"OS-EXT-STS:vm_state": "active",
			"OS-EXT-SRV-ATTR:instance_name": "instance-00000020",
			"OS-SRV-USG:launched_at": "2016-07-25T20:33:43.000000",
			"flavor": {
				"id": "8f8e2095-3b20-4617-b336-eadd5d698b64",
				"links": [
					{
						"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/flavors/8f8e2095-3b20-4617-b336-eadd5d698b64",
						"rel": "bookmark"
					}
				]
			},
			"id": "ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
			"security_groups": [
				{
					"name": "default"
				}
			],
			"scheduler_hints": {},
			"user_id": "0159cfde981a4f0cab47d59bd1a68601",
			"OS-DCF:diskConfig": "MANUAL",
			"accessIPv4": "",
			"accessIPv6": "",
			"progress": 0,
			"OS-EXT-STS:power_state": 1,
			"OS-EXT-AZ:availability_zone": "nova",
			"config_drive": "True",
			"status": "ACTIVE",
			"updated": "2016-07-25T20:33:44Z",
			"hostId": "92147c74df3beb53644ef5e569f6d0e5ab8d6f4a364486fc97ffe738",
			"OS-EXT-SRV-ATTR:host": "compute-0-1.domain.tld",
			"OS-SRV-USG:terminated_at": null,
			"key_name": null,
			"OS-EXT-SRV-ATTR:hypervisor_hostname": "compute-0-1.domain.tld",
			"name": "atlas",
			"created": "2016-07-25T20:33:10Z",
			"tenant_id": "6ee644b8e74b40448b768cba69249053",
			"os-extended-volumes:volumes_attached": [],
			"metadata": {
				"evacuation_policy": "NoEvacuation"
			}
		}
	]
}
`


const ServerErpResponseBody = `
{
	"server": {
		"OS-EXT-STS:task_state": null,
		"addresses": {
			"private": [
				{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:80:35:77",
					"version": 4,
					"addr": "10.33.160.11",
					"OS-EXT-IPS:type": "fixed"
				}
			],
		},
		"links": [
			{
				"href": "http://192.168.2.21:8774/v2/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
				"rel": "self"
			},
			{
				"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
				"rel": "bookmark"
			}
		],
		"image": {
			"id": "d74d06f2-61f5-4d02-b568-9d5002cf27f4",
			"links": [
				{
					"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/images/d74d06f2-61f5-4d02-b568-9d5002cf27f4",
					"rel": "bookmark"
				}
			]
		},
		"OS-EXT-STS:vm_state": "active",
		"OS-EXT-SRV-ATTR:instance_name": "instance-00000020",
		"OS-SRV-USG:launched_at": "2016-07-25T20:33:43.000000",
		"flavor": {
			"id": "8f8e2095-3b20-4617-b336-eadd5d698b64",
			"links": [
				{
					"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/flavors/8f8e2095-3b20-4617-b336-eadd5d698b64",
					"rel": "bookmark"
				}
			]
		},
		"id": "ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
		"security_groups": [
			{
				"name": "default"
			}
		],
		"scheduler_hints": {},
		"user_id": "0159cfde981a4f0cab47d59bd1a68601",
		"OS-DCF:diskConfig": "MANUAL",
		"accessIPv4": "",
		"accessIPv6": "",
		"progress": 0,
		"OS-EXT-STS:power_state": 1,
		"OS-EXT-AZ:availability_zone": "nova",
		"config_drive": "True",
		"status": "ACTIVE",
		"updated": "2016-07-25T20:33:44Z",
		"hostId": "92147c74df3beb53644ef5e569f6d0e5ab8d6f4a364486fc97ffe738",
		"OS-EXT-SRV-ATTR:host": "compute-0-1.domain.tld",
		"OS-SRV-USG:terminated_at": null,
		"key_name": null,
		"OS-EXT-SRV-ATTR:hypervisor_hostname": "compute-0-1.domain.tld",
		"name": "atlas",
		"created": "2016-07-25T20:33:10Z",
		"tenant_id": "6ee644b8e74b40448b768cba69249053",
		"os-extended-volumes:volumes_attached": [],
		"metadata": {
			"evacuation_policy": "NoEvacuation"
		}
	}
`

var (
	// Taken from gophercloud/openstack/v2/servers/fixtutes
	// ServerHerp is a Server struct that should correspond to the first result in ServerListBody.
	ServerHerp = servers.Server{
		Status:  "ACTIVE",
		Updated: "2014-09-25T13:10:10Z",
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]interface{}{
			"private": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:7c:1b:2b",
					"version":                 float64(4),
					"addr":                    "10.0.0.32",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []interface{}{
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "self",
			},
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
				"rel":  "bookmark",
			},
		},
		Image: map[string]interface{}{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]interface{}{
			"id": "1",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "ef079b0c-e610-4dfb-b1aa-b49f07ac48e5",
		UserID:   "9349aff8be7545ac9d2f1d00999a23cd",
		Name:     "herp",
		Created:  "2014-09-25T13:10:02Z",
		TenantID: "fcad67a6189847c4aecfa3c81a05783b",
		Metadata: map[string]interface{}{
			"evacuation_policy": "Evacuation",
			"evacuation_range": 100,
		},
		SecurityGroups: []map[string]interface{}{
			map[string]interface{}{
				"name": "default",
			},
		},
	}

	// ServerDerp is a Server struct that should correspond to the second server in ServerListBody.
	ServerDerp = servers.Server{
		Status:  "ACTIVE",
		Updated: "2014-09-25T13:04:49Z",
		HostID:  "29d3c8c896a45aa4c34e52247875d7fefc3d94bbcc9f622b5d204362",
		Addresses: map[string]interface{}{
			"private": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:9e:89:be",
					"version":                 float64(4),
					"addr":                    "10.0.0.31",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []interface{}{
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/v2/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "self",
			},
			map[string]interface{}{
				"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/servers/9e5476bd-a4ec-4653-93d6-72c93aa682ba",
				"rel":  "bookmark",
			},
		},
		Image: map[string]interface{}{
			"id": "f90f6034-2570-4974-8351-6b49732ef2eb",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/images/f90f6034-2570-4974-8351-6b49732ef2eb",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]interface{}{
			"id": "1",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://104.130.131.164:8774/fcad67a6189847c4aecfa3c81a05783b/flavors/1",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "9e5476bd-a4ec-4653-93d6-72c93aa682ba",
		UserID:   "9349aff8be7545ac9d2f1d00999a23cd",
		Name:     "derp",
		Created:  "2014-09-25T13:04:41Z",
		TenantID: "fcad67a6189847c4aecfa3c81a05783b",
		Metadata: map[string]interface{}{
			"evacuation_policy": "Evacuation",
			"evacuation_range": 0,
		},
		SecurityGroups: []map[string]interface{}{
			map[string]interface{}{
				"name": "default",
			},
		},
	}
	ServerErp = servers.Server{
		Status:  "ACTIVE",
		Updated: "2016-07-25T20:33:44Z",
		HostID:  "92147c74df3beb53644ef5e569f6d0e5ab8d6f4a364486fc97ffe738",
		Addresses: map[string]interface{}{
			"private": []interface{}{
				map[string]interface{}{
					"OS-EXT-IPS-MAC:mac_addr": "fa:16:3e:80:35:77",
					"version":                 float64(4),
					"addr":                    "10.33.160.11",
					"OS-EXT-IPS:type":         "fixed",
				},
			},
		},
		Links: []interface{}{
			map[string]interface{}{
				"href": "http://192.168.2.21:8774/v2/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
				"rel":  "self",
			},
			map[string]interface{}{
				"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/servers/ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
				"rel":  "bookmark",
			},
		},
		Image: map[string]interface{}{
			"id": "d74d06f2-61f5-4d02-b568-9d5002cf27f4",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/images/d74d06f2-61f5-4d02-b568-9d5002cf27f4",
					"rel":  "bookmark",
				},
			},
		},
		Flavor: map[string]interface{}{
			"id": "8f8e2095-3b20-4617-b336-eadd5d698b64",
			"links": []interface{}{
				map[string]interface{}{
					"href": "http://192.168.2.21:8774/6ee644b8e74b40448b768cba69249053/flavors/8f8e2095-3b20-4617-b336-eadd5d698b64",
					"rel":  "bookmark",
				},
			},
		},
		ID:       "ce0c725c-6d0b-4b4e-bd67-3129e9fed461",
		UserID:   "0159cfde981a4f0cab47d59bd1a68601",
		Name:     "erp",
		Created:  "2016-07-25T20:33:10Z",
		TenantID: "6ee644b8e74b40448b768cba69249053",
		Metadata: map[string]interface{}{
			"evacuation_policy": "NoEvacuation",
		},
		SecurityGroups: []map[string]interface{}{
			map[string]interface{}{
				"name": "default",
			},
		},
	}

	ServersList = []servers.Server{ServerErp, ServerHerp, ServerDerp}
	ServersListSortedByRangeExpected = []servers.Server{ServerDerp, ServerHerp, ServerErp}
	ServersListFilteredByEvacuationPolicyExpected = []servers.Server{ServerDerp, ServerHerp}
)

