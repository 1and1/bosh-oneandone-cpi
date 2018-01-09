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
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

var _ = Describe("HasDisk", func() {
	var (
		err   error
		found bool

		connector  *clientfakes.FakeConnector
		logger     boshlog.Logger
		diskFinder *diskfakes.FakeDiskFinder

		hasDisk HasDisk
	)

	BeforeEach(func() {
		diskFinder = &diskfakes.FakeDiskFinder{}
		installDiskFinderFactory(func(c client.Connector, l boshlog.Logger) disks.Finder {
			return diskFinder
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		hasDisk = NewHasDisk(connector, logger)
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		It("returns true if disk exists", func() {

			var identity sdk.Identity
			identity.Id = "fake-storage-id"
			identity.Name = "fake-name"

			diskFinder.FindStorageResult = &sdk.BlockStorage{
				Size:     20,
				Identity: identity,
			}


			found, err = hasDisk.Run("fake-storage-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(diskFinder.FindStorageID).To(Equal("fake-storage-id"))
		})

		It("returns false if disk ID does not exist", func() {

			var identity sdk.Identity
			identity.Id = "fake-storage-id"
			identity.Name = "fake-name"

			diskFinder.FindStorageResult = &sdk.BlockStorage{
				Size:     20,
				Identity: identity,
			}

			found, err = hasDisk.Run("fake-storage-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
		})

		It("returns an error if disk finder fails", func() {
			diskFinder.FindStorageError = errors.New("fake-find-vol-error")

			_, err = hasDisk.Run("fake-storage-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-find-vol-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
		})
	})
})
