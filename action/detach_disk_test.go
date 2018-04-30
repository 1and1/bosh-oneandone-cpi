package action

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	diskfakes "github.com/bosh-oneandone-cpi/oneandone/disks/fakes"
	vmfakes "github.com/bosh-oneandone-cpi/oneandone/vm/fakes"
	registryfakes "github.com/bosh-oneandone-cpi/registry/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	"github.com/bosh-oneandone-cpi/registry"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"

	"github.com/bosh-oneandone-cpi/oneandone/vm"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/disks"
)

var _ = Describe("DetachDisk", func() {
	var (
		err        error
		detacherVM *resource.Instance

		foundInstance *resource.Instance
		foundVolume   *sdk.BlockStorage
		expectedAgentSettings registry.AgentSettings

		registryClient *registryfakes.FakeClient
		connector      *clientfakes.FakeConnector
		logger         boshlog.Logger

		vmFinder *vmfakes.FakeVMFinder

		diskFinder       *diskfakes.FakeDiskFinder
		attacherDetacher *diskfakes.FakeAttacherDetacher

		detachDisk DetachDisk
	)
	BeforeEach(func() {
		vmFinder = &vmfakes.FakeVMFinder{}
		installVMFinderFactory(func(c client.Connector, l boshlog.Logger) vm.Finder {
			return vmFinder
		})

		diskFinder = &diskfakes.FakeDiskFinder{}
		installDiskFinderFactory(func(c client.Connector, l boshlog.Logger) disks.Finder {
			return diskFinder
		})

		attacherDetacher = &diskfakes.FakeAttacherDetacher{}
		installInstanceAttacherDetacherFactory(func(in *resource.Instance, c client.Connector, l boshlog.Logger) (disks.AttacherDetacher, error) {
			detacherVM = in
			return attacherDetacher, nil
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		registryClient = &registryfakes.FakeClient{}

		detachDisk = NewDetachDisk(connector, logger, registryClient)

		foundInstance = resource.NewInstance("fake-vm-ocid")
		vmFinder.FindInstanceResult = foundInstance

		diskFinder.FindStorageResult = foundVolume

		registryClient.FetchSettings = registry.AgentSettings{
			Disks: registry.DisksSettings{
				Persistent: map[string]registry.PersistentSettings{
					"fake-vol-ocid": {
						ID:   "fake-vol-ocid",
						Path: "/dev/fake-path",
					},
				},
			},
		}
		expectedAgentSettings = registry.AgentSettings{
			Disks: registry.DisksSettings{
				Persistent: map[string]registry.PersistentSettings{},
			},
		}
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		It("finds the vm", func() {
			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceID).To(Equal("fake-vm-ocid"))
		})

		It("creates a detacher for the found vm", func() {
			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(detacherVM).To(Equal(vmFinder.FindInstanceResult))
		})

		It("finds the disk", func() {
			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(diskFinder.FindStorageID).To(Equal("fake-vol-ocid"))
		})

		It("detaches the disk", func() {
			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(attacherDetacher.DetachVolumeCalled).To(BeTrue())
			Expect(attacherDetacher.DetachedVolume).To(Equal(foundVolume))
		})

		It("udates the registry", func() {
			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))

		})
		It("returns an error if vmfinder fails", func() {
			vmFinder.FindInstanceError = errors.New("fake-instance-finder-error")

			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("fake-instance-finder-error"))
			Expect(attacherDetacher.DetachVolumeCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if diskfinder fails", func() {
			diskFinder.FindStorageError = errors.New("fake-disk-finder-error")

			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-finder-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(attacherDetacher.DetachVolumeCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if detacher fails", func() {
			attacherDetacher.DetachmentError = errors.New("fake-attachment-error")

			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-attachment-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.DetachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.DetachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = detachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.DetachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
