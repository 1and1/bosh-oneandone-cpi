package integration

import (
	"fmt"

	"github.com/1and1/oneandone-cloudserver-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var machineName = "bosh director integration test"
var _ = Describe("VM", func() {
	It("executes the VM lifecycle", func() {
		var vmCID string

		By("creating a VM with wrong ssh public key passed to API")
		request := fmt.Sprintf(`{
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
			}`, existingStemcell, machineName, "WRONG SSH KEY", pnNetworkId)
		resp, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Error.Message).ToNot(BeEmpty())
		Expect(resp.Log).ToNot(BeEmpty())

		By("creating a VM with no Private network id")
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
			}`, existingStemcell, machineName, sshKey)
		resp, err = execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Error.Message).ToNot(BeEmpty())
		Expect(resp.Log).ToNot(BeEmpty())

		By("creating a VM with an invalid hardware flavor")
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "name": "%v",
				  "flavor": "1S",
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
		resp, err = execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Error.Message).ToNot(BeEmpty())
		Expect(resp.Log).ToNot(BeEmpty())

		By("creating a VM with no rsa key provided")
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "name": "%v",
				  "flavor": "S"
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
			}`, existingStemcell, machineName, pnNetworkId)
		resp, err = execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.Error.Message).ToNot(BeEmpty())
		Expect(resp.Log).ToNot(BeEmpty())

		By("creating a VM with custom hardware settings")
		request = fmt.Sprintf(`{
			  "method": "create_vm",
			  "arguments": [
				"agent",
				"%v",
				{
				  "name": "%v",
 				  "cores": 2,
                  "diskSize":30,
                  "ram":4,
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

		By("creating a VM")
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

		By("locating the VM")
		request = fmt.Sprintf(`{
			  "method": "has_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		exists := assertSucceedsWithResult(request).(bool)
		Expect(exists).To(Equal(true))

		assertValidVM(vmCID, func(instance *oneandone.Server) {
			Expect(instance.Name).To(ContainSubstring(machineName))
		})

		updatedName := "updatedfrombosh"
		request = fmt.Sprintf(`{
			  "method": "set_vm_metadata",
			  "arguments": [
				"%v",
				{"name":"%v"}
			  ]
			}`, vmCID, updatedName)
		assertSucceeds(request)
		assertValidVM(vmCID, func(instance *oneandone.Server) {
			Expect(instance.Name).To(ContainSubstring(updatedName))
		})

		By("deleting the VM")
		request = fmt.Sprintf(`{
			  "method": "delete_vm",
			  "arguments": ["%v"]
			}`, vmCID)
		assertSucceeds(request)

	})
})
