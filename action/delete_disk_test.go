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
)

var _ = Describe("DeleteDisk", func() {
	var (
		err error

		connector      *clientfakes.FakeConnector
		logger         boshlog.Logger
		diskTerminator *diskfakes.FakeDiskTerminator

		deleteDisk DeleteDisk
	)

	BeforeEach(func() {
		diskTerminator = &diskfakes.FakeDiskTerminator{}
		installDiskTerminatorFactory(func(c client.Connector, l boshlog.Logger) disks.Terminator {
			return diskTerminator
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)

		deleteDisk = NewDeleteDisk(connector, logger)
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		It("delegates to disk terminator", func() {
			_, err = deleteDisk.Run("fake-volume-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskTerminator.DeleteStorageCalled).To(BeTrue())
			Expect(diskTerminator.DeleteStorageID).To(Equal("fake-volume-ocid"))
		})

		It("returns an error if disk terminator returns an error", func() {
			diskTerminator.DeleteStorageError = errors.New("fake-delete-volume-error")

			_, err = deleteDisk.Run("fake-volume-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-delete-volume-error"))
			Expect(diskTerminator.DeleteStorageCalled).To(BeTrue())
		})
	})
})
