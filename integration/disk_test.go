package integration

//import (
//	"fmt"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//)

//var _ = Describe("Disk", func() {
//
//	It("executes the disk lifecycle", func() {
//		By("creating a disk")
//		var diskCID string
//		request := fmt.Sprintf(`{
//			  "method": "create_disk",
//			  "arguments": [20,{"datacenter":"908DC2072407C94C8054610AD5A53B8C"}]
//			}`)
//		diskCID = assertSucceedsWithResult(request).(string)
//
//		By("confirming a disk exists")
//		request = fmt.Sprintf(`{
//			  "method": "has_disk",
//			  "arguments": ["%v"]
//			}`, diskCID)
//		assertSucceeds(request)
//
//		By("creating a VM")
//		var vmCID string
//		request = fmt.Sprintf(`{
//			  "method": "create_vm",
//			  "arguments": [
//				"agent",
//				"%v",
//				{
//				  "Name": "boshtest",
//				  "flavor":"s",
//    "rsa_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDA3zMuWZBH4INpQ/XegqfJ498eQHyffuqG1Y/NSzUNAY2Ys6xggQijATbAIF6GuCSfxfKzlmR2J6HmhxzYqhhmNZk7vIyfQ6TWYyoNnDNCcysFy+GOf65WY9KrUo3JH8SJtytn5Hs6ZT7GT22pc2P3ThiHu3e+y3N3nhjglgLXnKpzQbS2xzeA8AUUfC7TfIIpE1JP+856BXH8b3a9rbt6xQQEQFdTqOAZgi93+wp9jHiiq8NsTWSi1W4t0+r0Y9XM0t1DGjpKUTm7Q+JxItnx6IpBtl6j8ThKo9J5NIlmAEKhDQYdecsDnjU3WhuENRAb/jfP2+Ci7i7bANVxa/Tj ali@ali-G751JT"
//				},
//				{
//				  "default": {
//					"type": "dynamic",
//					"cloud_properties": {
//					  "open-ports": [
//							{
//								"port-from":22,
//								"port-to":22,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":80,
//								"port-to":80,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":443,
//								"port-to":443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8443,
//								"port-to":8443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8447,
//								"port-to":8447,
//								"source":"0.0.0.0"
//							}
//						]
//					}
//				  }
//				},
//				{
//				  "bosh": {
//					  "group_name": "1and1 test",
//					  "groups": ["micro-1and1", "dummy", "dummy", "micro-1and1-dummy", "dummy-dummy"]
//				  }
//				}
//			  ]
//			}`, existingStemcell)
//		vmCID = assertSucceedsWithResult(request).(string)
//
//		By("attaching the disk")
//		request = fmt.Sprintf(`{
//			  "method": "attach_disk",
//			  "arguments": ["%v", "%v"]
//			}`, vmCID, diskCID)
//		assertSucceeds(request)
//
//		By("confirming the attachment of a disk")
//		request = fmt.Sprintf(`{
//			  "method": "get_disks",
//			  "arguments": ["%v"]
//			}`, vmCID)
//		disks := toStringArray(assertSucceedsWithResult(request).([]interface{}))
//		Expect(disks).To(ContainElement(diskCID))
//
//		By("detaching and deleting a disk")
//		request = fmt.Sprintf(`{
//			  "method": "detach_disk",
//			  "arguments": ["%v", "%v"]
//			}`, vmCID, diskCID)
//		assertSucceeds(request)
//
//		request = fmt.Sprintf(`{
//			  "method": "delete_disk",
//			  "arguments": ["%v"]
//			}`, diskCID)
//		assertSucceeds(request)
//
//		By("deleting the VM")
//		request = fmt.Sprintf(`{
//			  "method": "delete_vm",
//			  "arguments": ["%v"]
//			}`, vmCID)
//		assertSucceeds(request)
//	})
//})
