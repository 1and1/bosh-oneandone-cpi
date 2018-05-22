package action

import (
	"errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/disks"
	diskfakes "github.com/bosh-oneandone-cpi/oneandone/disks/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	vmfakes "github.com/bosh-oneandone-cpi/oneandone/vm/fakes"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
)

var _ = Describe("CreateDisk", func() {
	var (
		err        error
		diskCID    DiskCID
		cloudProps DiskCloudProperties

		connector   *clientfakes.FakeConnector
		logger      boshlog.Logger
		diskCreator *diskfakes.FakeDiskCreator
		vmFinder    *vmfakes.FakeVMFinder

		createDisk CreateDisk
	)

	BeforeEach(func() {

		vmFinder = &vmfakes.FakeVMFinder{}
		installVMFinderFactory(func(c client.Connector, l boshlog.Logger) vm.Finder {
			return vmFinder
		})

		diskCreator = &diskfakes.FakeDiskCreator{}
		installDiskCreatorFactory(func(c client.Connector, l boshlog.Logger) disks.Creator {
			return diskCreator
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		createDisk = NewCreateDisk(connector, logger)
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		Context("when vmCID instance is found", func() {

			BeforeEach(func() {
				vmFinder.FindInstanceResult = resource.NewInstance("fake-vm-ocid")
				cloudProps = DiskCloudProperties{}
			})

			It("creates the disk", func() {
				var identity sdk.Identity
				identity.Id = "fake-storage-id"
				identity.Name = "fake-name"

				diskCreator.CreateStorageResult = &sdk.BlockStorage{
					Size:     20,
					Identity: identity,
				}

				diskCID, err = createDisk.Run(20, cloudProps)

				Expect(err).NotTo(HaveOccurred())
				Expect(diskCreator.CreateStorageCalled).To(BeTrue())
				Expect(diskCID).To(Equal(DiskCID("fake-storage-id")))
			})
			It("returns an error if disk creator fails", func() {
				diskCreator.CreateStorageError = errors.New("fake-create-storage-error")

				_, err = createDisk.Run(20, cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-create-storage-error"))
			})

			It("rounds up the requested disk size to the minimum supported size", func() {
				var identity sdk.Identity
				identity.Id = "fake-storage-id"
				identity.Name = "fake-name"

				diskCreator.CreateStorageResult = &sdk.BlockStorage{
					Size:     10,
					Identity: identity,
				}
				_, err = createDisk.Run(10, cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(diskCreator.CreateStorageCalled).To(BeTrue())
				Expect(diskCreator.CreateStorageSize >= minVolumeSize).To(BeTrue())
			})

		})
	})
})
