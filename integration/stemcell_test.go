package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	It("executes the stemcell lifecycle with an image-id", func() {

		By("finding a stemcell by image-id")
		request := fmt.Sprintf(`{
         "method": "create_stemcell",
         "arguments": ["",{
           "name": "bosh-oneandone-xen-ubuntu-trusty-go_agent-raw",
           "version": "3541.21",
           "infrastructure": "oneandone",
           "image-id": "%s"
         }]
       }`, imageId)

		response, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Error).To(BeNil())
		Expect(response.Result).To(Not(BeNil()))
	})
})
