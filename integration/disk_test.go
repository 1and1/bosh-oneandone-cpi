package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Disk", func() {

	It("executes the disk lifecycle", func() {
		By("creating a disk")
		var diskCID string
		request := fmt.Sprintf(`{
			  "method": "create_disk",
			  "arguments": [20,{"datacenter":"908DC2072407C94C8054610AD5A53B8C"}]
			}`)
		diskCID = assertSucceedsWithResult(request).(string)

		By("confirming a disk exists")
		request = fmt.Sprintf(`{
			  "method": "has_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		By("creating a VM")
		var vmCID string
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "name": "%v",
				  "flavor": "S",
                  "keypair":"/root/.ssh",
				  "rsa_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD0IdDhD9pzUYUBEmD8sbUcisr6kTh8a4eOmdN5DI3WkJIO3NhVWWHMZfBMApJHTVpgKOcrmArYZpft08QPOiRb2Om/0nTQvXLAjo/ra0lUYrHQw8WZW88Itzf1mSHN3dlsc+YoJPSFeRksqpntWnL/TwLyuJQ51qxIew+RTitayDdRtR+Qhn1qw/yxtH4Mt+nFJMu4OORBCR3CdrcAHUmmBOZ3eOr2WHWuTHVDrSuqgqc7ndnABWwQOs37fKsL38tEZC0oKbHM34alizSmXjszzIMMM3HMoDyS4cDBdS8uoNaSU1/fMZj3BkTQST+UwJtLZN+3X/ClKJztz9ijwYMR root@ali-G751JT"

				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
                     "private-network-id":"D522A56E643EED2479F2B73810DAF5F3",
					 "open-ports": [
							{
								"port-from":22,
								"port-to":22,
								"source":"0.0.0.0"

							},
							{
								"port-from":80,
								"port-to":80,
								"source":"0.0.0.0"
							}
						]
					}
				  }
				},["diskid"],
				{
				  "bosh": {
					  "group_name": "1and1 test",
				  "groups": ["micro-1and1", "dummy", "dummy", "micro-1and1-dummy", "dummy-dummy"]
				  }
				}
			  ]
			}`, existingStemcell, machineName)
		vmCID = assertSucceedsWithResult(request).(string)

		By("attaching the disk")
		request = fmt.Sprintf(`{
			  "method": "attach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

		By("confirming the attachment of a disk")
		request = fmt.Sprintf(`{
			  "method": "get_disks",
			  "arguments": ["%v"]
			}`, vmCID)
		disks := toStringArray(assertSucceedsWithResult(request).([]interface{}))
		Expect(disks).To(ContainElement(diskCID))

		By("detaching and deleting a disk")
		request = fmt.Sprintf(`{
			  "method": "detach_disk",
			  "arguments": ["%v", "%v"]
			}`, vmCID, diskCID)
		assertSucceeds(request)

		request = fmt.Sprintf(`{
			  "method": "delete_disk",
			  "arguments": ["%v"]
			}`, diskCID)
		assertSucceeds(request)

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)
	})
})
