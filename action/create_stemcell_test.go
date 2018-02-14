package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"

	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	stemcellfakes "github.com/bosh-oneandone-cpi/oneandone/stemcell/fakes"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/stemcell"
)

var _ = Describe("CreateStemcell", func() {
	var (
		cloudProps StemcellCloudProperties

		connector      *clientfakes.FakeConnector
		logger         boshlog.Logger
		creator        *stemcellfakes.FakeCreator
		finder         *stemcellfakes.FakeFinder
		createStemcell CreateStemcell
	)

	BeforeEach(func() {

		finder = &stemcellfakes.FakeFinder{}
		creator = &stemcellfakes.FakeCreator{}
		installStemcellCreatorFactory(func(c client.Connector, l boshlog.Logger) stemcell.Creator {
			return creator
		})
		installStemcellFinderFactory(func(c client.Connector, l boshlog.Logger) stemcell.Finder {
			return finder
		})
		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		createStemcell = NewCreateStemcell(connector, logger)
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		Context("When called with image-ocid property set", func() {
			BeforeEach(func() {
				cloudProps = StemcellCloudProperties{
					ImageID: "fake-image-ocid",
				}
			})
		})
	})
})
