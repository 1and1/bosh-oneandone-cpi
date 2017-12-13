package action

import (
	"errors"
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	"github.com/bosh-oneandone-cpi/registry"

	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	vmfakes "github.com/bosh-oneandone-cpi/oneandone/vm/fakes"
	registryfakes "github.com/bosh-oneandone-cpi/registry/fakes"
	fakeuuid "github.com/cloudfoundry/bosh-utils/uuid/fakes"
)

var _ = Describe("CreateVM", func() {
	var (
		connector      *clientfakes.FakeConnector
		logger         boshlog.Logger
		registryClient *registryfakes.FakeClient
		uuidGen        *fakeuuid.FakeGenerator
		creator        *vmfakes.FakeVMCreator

		err             error
		cloudProps      VMCloudProperties
		env             Environment
		registryOptions registry.ClientOptions
		agentOptions    registry.AgentOptions

		expectedAgentSettings registry.AgentSettings
		vmCID                 VMCID

		createVM CreateVM
	)

	BeforeEach(func() {
		installVMCreatorFactory(func(c client.Connector, l boshlog.Logger) vm.Creator {
			return creator
		})

		connector = &clientfakes.FakeConnector{}
		uuidGen = &fakeuuid.FakeGenerator{}
		uuidGen.GeneratedUUID = "fake-uuid"
		logger = boshlog.NewLogger(boshlog.LevelNone)
		registryClient = &registryfakes.FakeClient{}
		registryOptions = registry.ClientOptions{
			Protocol: "http",
			Host:     "fake-registry-host",
			Port:     25777,
			Username: "fake-registry-username",
			Password: "fake-registry-password",
		}
		agentOptions = registry.AgentOptions{
			Mbus: "http://fake-mbus",
			Blobstore: registry.BlobstoreOptions{
				Provider: "fake-blobstore-type",
				Options:  nil,
			},
		}
		connector.AgentOptionsResult = agentOptions
		createVM = NewCreateVM(connector, logger, registryClient, uuidGen)
	})
	AfterEach(func() { resetAllFactories() })

	Describe("Run", func() {
		BeforeEach(func() {

			cloudProps = VMCloudProperties{
				Cores:      1,
				Ram:        4,
				DiskSize:   30,
				Name:       "test",
				Datacenter: "fake-dc",
			}
		})
		Context("when no errors in vm creation", func() {
			BeforeEach(func() {
				creator = &vmfakes.FakeVMCreator{
					CreateInstanceResult: resource.NewInstance("fake-ocid"),
					CreateInstanceError:  nil,
				}
			})
			It("creates the vm", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, nil, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(creator.CreateInstanceCalled).To(BeTrue())
				Expect(vmCID).To(Equal(VMCID("fake-ocid")))
			})
			It("uses uuid as part of vm name", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, nil, env)
				Expect(creator.Configuration.Name).To(ContainSubstring(uuidGen.GeneratedUUID))
			})
			It("updates the registry", func() {
				expectedAgentSettings = registry.AgentSettings{
					AgentID: "fake-agent-id",
					Blobstore: registry.BlobstoreSettings{
						Provider: "fake-blobstore-type",
						Options:  nil,
					},
					Disks: registry.DisksSettings{
						System:     "/dev/sda",
						Ephemeral:  "/dev/sda",
						Persistent: map[string]registry.PersistentSettings{},
					},
					Mbus: "http://fake-mbus",
					VM: registry.VMSettings{
						Name: fmt.Sprintf("bosh-%s", uuidGen.GeneratedUUID),
					},
				}
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, nil, env)
				Expect(err).NotTo(HaveOccurred())
				Expect(registryClient.UpdateCalled).To(BeTrue())
				Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
				Expect(vmCID).To(Equal(VMCID("fake-ocid")))
			})
		})
		Context("when vm creation fails", func() {
			BeforeEach(func() {
				creator = &vmfakes.FakeVMCreator{
					CreateInstanceResult: nil,
					CreateInstanceError:  errors.New("fake-launchinstance-error"),
				}
			})
			It("propagates the error", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, nil, env)
				Expect(err).To(HaveOccurred())
				Expect(creator.CreateInstanceCalled).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("fake-launchinstance-error"))
				Expect(vmCID).To(Equal(VMCID("")))
			})
			It("doesn't update registry", func() {
				vmCID, err = createVM.Run("fake-agent-id", "fake-stemcell-id", cloudProps, nil, env)
				Expect(err).To(HaveOccurred())
				Expect(creator.CreateInstanceCalled).To(BeTrue())
				Expect(registryClient.UpdateCalled).To(BeFalse())
				Expect(vmCID).To(Equal(VMCID("")))
			})
		})

	})

})
