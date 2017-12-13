package action

import (
	"errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	stemcellfakes "github.com/bosh-oneandone-cpi/oneandone/stemcell/fakes"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/stemcell"
)

var _ = Describe("CreateStemcell", func() {
	var (
		err        error
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
			It("returns an upgrade stemcell error", func() {
				_, err = createStemcell.Run("", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(creator.CreateStemcellCalled).To(BeFalse())
				Expect(finder.FindStemcellCalled).To(BeFalse())
			})

		})
		Context("When called with image source URL property set", func() {
			BeforeEach(func() {
				cloudProps = StemcellCloudProperties{
					Name:           "fake-image-name",
					ImageSourceURL: "fake-image-source-url",
				}
			})
			It("delegates to StemCell Creator", func() {
				_, err = createStemcell.Run("", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(finder.FindStemcellCalled).To(BeFalse())
				Expect(creator.CreateStemcellCalled).To(BeTrue())
			})
			It("passes the image source url to creator", func() {
				_, err = createStemcell.Run("", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(creator.CreateStemcellCalledWithURL).To(Equal("fake-image-source-url"))
			})
			It("passes concatenated name and version to creator", func() {
				_, err = createStemcell.Run("", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(creator.CreateStemcellCalledWithImageName).To(ContainSubstring("fake-image-name"))

			})

			It("returns error if stemcell creator fails ", func() {
				creator.CreateStemcellError = errors.New("fake-create-image-error")
				_, err = createStemcell.Run("", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(creator.CreateStemcellCalled).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("fake-create-image-error"))
			})

		})
		Context("When called with ImageSourceURL not set", func() {
			It("returns an error", func() {
				_, err = createStemcell.Run("", StemcellCloudProperties{})
				Expect(err).To(HaveOccurred())
				Expect(creator.CreateStemcellCalled).To(BeFalse())
				Expect(finder.FindStemcellCalled).To(BeFalse())
			})

		})
	})

})
