package integration

import (
	"fmt"

	"github.com/oneandone/oneandone-cloudserver-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {
	//request := fmt.Sprintf(`{
	//		  "method": "create_vm",
	//		  "arguments": [
	//			"agent",
	//			"%v",
	//			{
	//			  "name": "boshtest",
	//			  "flavor": "S",
	//			  "rsa_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD0IdDhD9pzUYUBEmD8sbUcisr6kTh8a4eOmdN5DI3WkJIO3NhVWWHMZfBMApJHTVpgKOcrmArYZpft08QPOiRb2Om/0nTQvXLAjo/ra0lUYrHQw8WZW88Itzf1mSHN3dlsc+YoJPSFeRksqpntWnL/TwLyuJQ51qxIew+RTitayDdRtR+Qhn1qw/yxtH4Mt+nFJMu4OORBCR3CdrcAHUmmBOZ3eOr2WHWuTHVDrSuqgqc7ndnABWwQOs37fKsL38tEZC0oKbHM34alizSmXjszzIMMM3HMoDyS4cDBdS8uoNaSU1/fMZj3BkTQST+UwJtLZN+3X/ClKJztz9ijwYMR root@ali-G751JT"
	//
	//			},
	//			{
	//			  "default": {
	//				"type": "dynamic",
	//				"cloud_properties": {
	//				  "open-ports": [
	//						{
	//							"port-from":22,
	//							"port-to":22,
	//							"source":"0.0.0.0"
	//
	//						},
	//						{
	//							"port-from":80,
	//							"port-to":80,
	//							"source":"0.0.0.0"
	//
	//						},
	//						{
	//							"port-from":443,
	//							"port-to":443,
	//							"source":"0.0.0.0"
	//
	//						},
	//						{
	//							"port-from":8443,
	//							"port-to":8443,
	//							"source":"0.0.0.0"
	//
	//						},
	//						{
	//							"port-from":8447,
	//							"port-to":8447,
	//							"source":"0.0.0.0"
	//						}
	//					]
	//				}
	//			  }
	//			}
	//		  ]
	//		}`, existingStemcell)
	//It("creates a VM with an invalid configuration and receives an error message with logs", func() {
	//	resp, err := execCPI(request)
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(resp.Error.Message).ToNot(BeEmpty())
	//})

	It("executes the VM lifecycle", func() {
		var vmCID string
		By("creating a VM")
		request := fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "name": "boshtest",
				  "flavor": "S",
				  "rsa_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQD0IdDhD9pzUYUBEmD8sbUcisr6kTh8a4eOmdN5DI3WkJIO3NhVWWHMZfBMApJHTVpgKOcrmArYZpft08QPOiRb2Om/0nTQvXLAjo/ra0lUYrHQw8WZW88Itzf1mSHN3dlsc+YoJPSFeRksqpntWnL/TwLyuJQ51qxIew+RTitayDdRtR+Qhn1qw/yxtH4Mt+nFJMu4OORBCR3CdrcAHUmmBOZ3eOr2WHWuTHVDrSuqgqc7ndnABWwQOs37fKsL38tEZC0oKbHM34alizSmXjszzIMMM3HMoDyS4cDBdS8uoNaSU1/fMZj3BkTQST+UwJtLZN+3X/ClKJztz9ijwYMR root@ali-G751JT"

				},
				{
				  "default": {
					"type": "dynamic",
					"cloud_properties": {
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
				},
				{
				  "bosh": {
					  "group_name": "1and1 test",
				  "groups": ["micro-1and1", "dummy", "dummy", "micro-1and1-dummy", "dummy-dummy"]
				  }
				}
			  ]
			}`, existingStemcell)
		vmCID = assertSucceedsWithResult(request).(string)

		By("locating the VM")
		request = fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		exists := assertSucceedsWithResult(request).(bool)
		Expect(exists).To(Equal(true))

		expectedName := "boshtest"
		assertValidVM(vmCID, func(instance *oneandone.Server) {
			// Labels should be an exact match
			Expect(instance.Name).To(BeEquivalentTo(expectedName))
		})

		updatedName := "updatedfrombosh"
		request = fmt.Sprintf(`{
			  "method": "set_vm_metadata",
			  "arguments": [
				"%v",
				%v
			  ]
			}`, vmCID, updatedName)
		assertSucceeds(request)
		assertValidVM(vmCID, func(instance *oneandone.Server) {
			// Labels should be an exact match
			Expect(instance.Name).To(BeEquivalentTo(expectedName))
		})

		//By("rebooting the VM")
		//request = fmt.Sprintf(`{
		//	  "method": "reboot_vm",
		//	  "arguments": ["%v"]
		//	}`, vmCID)
		//assertSucceeds(request)

		//By("deleting the VM")
		//request = fmt.Sprintf(`{
		//	  "method": "delete_vm",
		//	  "arguments": ["%v"]
		//	}`, vmCID)
		//assertSucceeds(request)

	})
})
