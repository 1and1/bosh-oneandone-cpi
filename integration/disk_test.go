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
			  "arguments": [20,{"datacenter":"%v"}]
			}`, datacenterId)
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
               	  "rsa_key": "%v"
				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
                    "private-network-id":"%v",
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

							},
							{
								"port-from":443,
								"port-to":443,
								"source":"0.0.0.0"

							},
							{
								"port-from":8443,
								"port-to":8443,
								"source":"0.0.0.0"

							},
							{
								"port-from":8447,
								"port-to":8447,
								"source":"0.0.0.0"
							}
						]
					}
				  }
				},["diskid"],
				{
				  "bosh": {
				  }
				}
			  ]
			}`, existingStemcell, machineName, sshKey, pnNetworkId)
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
