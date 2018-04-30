package action

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	diskfakes "github.com/bosh-oneandone-cpi/oneandone/disks/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/disks"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
)

var _ = Describe("GetDisks", func() {
	var (
		err             error
		storages         []string

		connector  *clientfakes.FakeConnector
		logger     boshlog.Logger
		diskFinder *diskfakes.FakeDiskFinder

		getDisks GetDisks
	)

	BeforeEach(func() {
		diskFinder = &diskfakes.FakeDiskFinder{}
		installDiskFinderFactory(func(c client.Connector, l boshlog.Logger) disks.Finder {
			return diskFinder
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		getDisks = NewGetDisks(connector, logger)
	})

	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		Context("when there are attached disk", func() {

			BeforeEach(func() {

				var storages []sdk.BlockStorage
				var identity sdk.Identity
				identity.Id = "fake-storage-id"
				identity.Name = "fake-name"

				storages= []sdk.BlockStorage{
					sdk.BlockStorage{
						Size:20,
						Identity:identity,
					},
					sdk.BlockStorage{
						Size:40,
						Identity:identity,
					},
				}

				diskFinder.FindAllAttachedResult = storages
			})

			It("returns the list of attached disks", func() {
				storages, err = getDisks.Run("fake-vm-ocid")
				Expect(err).NotTo(HaveOccurred())
				Expect(diskFinder.FindAllAttachedStoragesCalled).To(BeTrue())
				Expect(storages).To(Equal([]string{"fake-storage-id", "fake-storage-id"}))
			})
		})

		Context("when there are not any attached disk", func() {
			It("returns an empty array", func() {
				storages, err = getDisks.Run("fake-vm-ocid")
				Expect(err).NotTo(HaveOccurred())
				Expect(diskFinder.FindAllAttachedStoragesCalled).To(BeTrue())
				Expect(storages).To(BeEmpty())
			})
		})

		It("returns an error if finder fails", func() {
			diskFinder.FindAllAttachedError = errors.New("fake-volfinder-error")

			_, err = getDisks.Run("fake-vm-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-volfinder-error"))
			Expect(diskFinder.FindAllAttachedStoragesCalled).To(BeTrue())
		})
	})
})
