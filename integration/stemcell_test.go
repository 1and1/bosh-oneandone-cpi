package integration

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stemcell", func() {
	It("executes the stemcell lifecycle with an image-id", func() {
		var imageId="BC7CC2DF4258653C41FFCBD18C4D7E9B"

		By("finding a stemcell by image-id")
		request := fmt.Sprintf(`{
         "method": "create_stemcell",
         "arguments": [{
           "name": "bosh-oneandone-kvm-ubuntu-trusty",
           "version": "3215",
           "infrastructure": "1&1",
           "image-id": "%s"
         }]
       }`, imageId)

		response, err := execCPI(request)
		Expect(err).ToNot(HaveOccurred())
		Expect(response.Error).To(BeNil())
		Expect(response.Result).To(Not(BeNil()))
	})
})
