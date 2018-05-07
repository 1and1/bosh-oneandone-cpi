package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"fmt"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
	"strings"
)

var imageId string
var pnNetworkId string
var privateNetworkName = "BOSH integration PN test"
var datacenterId string

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	initAPI()
	// Clean any leftovers
	cleanVMs()

	//find US datacenter
	datacetners, err := oaoClient.Client().ListDatacenters()
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	for _, dc := range datacetners {
		if dc.CountryCode == "US" {
			datacenterId = dc.Id
		}
	}

	images, err := oaoClient.Client().ListImages(1, 20, "", "bosh", "")
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}

	if len(images) > 0 {
		imageId = images[0].Id
	} else {
		Fail("BOSH image not found")
	}

	request := fmt.Sprintf(`{
			  "method": "create_stemcell",
         "arguments": ["",
			{
           "name": "bosh-oneandone-kvm-ubuntu-trusty",
           "version": "3215",
           "infrastructure": "1&1",
           "image-id": "%s"
				}]
			}`, imageId)
	stemcell := assertSucceedsWithResult(request).(string)

	pnRequest := sdk.PrivateNetworkRequest{
		Name: privateNetworkName,
	}
	_, pn, err := oaoClient.Client().CreatePrivateNetwork(&pnRequest)
	if err != nil {
		Expect(err).ToNot(HaveOccurred())
	}
	oaoClient.Client().WaitForState(pn, "ACTIVE", 10, 60)
	pnNetworkId = pn.Id

	return []byte(stemcell)
}, func(data []byte) {
	// Ensure stemcell was initialized
	existingStemcell = string(data)
	Expect(existingStemcell).ToNot(BeEmpty())

})

var _ = SynchronizedAfterSuite(func() {
	// Clean any leftovers
	cleanVMs()
}, func() {
})

func cleanVMs() {
	//delete dangling servers from previous tests
	serversToDelete, err := oaoClient.Client().ListServers(1, 20, "", machineName, "")
	if err == nil {

		for _, vm := range serversToDelete {
			GinkgoWriter.Write([]byte(fmt.Sprintf("Deleting VM %v\n", vm.Name)))
			if strings.Contains(vm.Name, machineName) && vm.Status.State == "POWERED_ON" {
				//remove firewall policies
				if len(vm.Ips) > 0 {
					for _, ip := range vm.Ips {
						fws, err := oaoClient.Client().GetServerIpFirewallPolicy(vm.Id, ip.Id)
						if err != nil {
							continue
						}
						delFw, err := oaoClient.Client().DeleteFirewallPolicy(fws.Id)
						if err != nil {
							continue
						}
						oaoClient.Client().WaitUntilDeleted(delFw)
					}

				}
				del, err := oaoClient.Client().DeleteServer(vm.Id, false)
				oaoClient.Client().WaitUntilDeleted(del)
				Expect(err).ToNot(HaveOccurred())
			}
		}
	}

	//delete dangling private networks from previous tests
	pns, err := oaoClient.Client().ListPrivateNetworks(1, 20, "", "bosh", "")
	if err == nil {

		for _, pn := range pns {
			GinkgoWriter.Write([]byte(fmt.Sprintf("Deleting Private Network %v\n", pn.Name)))
			if strings.Contains(pn.Name, privateNetworkName) && pn.State == "ACTIVE" {
				del, err := oaoClient.Client().DeletePrivateNetwork(pn.Id)
				oaoClient.Client().WaitUntilDeleted(del)
				Expect(err).ToNot(HaveOccurred())
			}
		}
	}

}
