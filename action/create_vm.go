package action

import (
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"time"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	"github.com/bosh-oneandone-cpi/registry"
)

// CreateVM action handles the create_vm request
type CreateVM struct {
	connector client.Connector
	logger    boshlog.Logger
	registry  registry.Client
	uuidGen   boshuuid.Generator
}

const logTag = "createVM"

// NewCreateVM creates a CreateVM instance
func NewCreateVM(c client.Connector, l boshlog.Logger, r registry.Client, u boshuuid.Generator) CreateVM {
	return CreateVM{connector: c, logger: l, registry: r, uuidGen: u}
}

func (cv CreateVM) Run(agentID string, stemcellCID StemcellCID, networks Networks, cloudProps VMCloudProperties, env Environment) (VMCID, error) {

	agentNetworks := networks.AsRegistryNetworks()
	// Create the VM
	name := cv.vmName()
	creator := newVMCreator(cv.connector, cv.logger)

	icfg := vm.InstanceConfiguration{
		Name:     name,
		ImageId:  string(stemcellCID),
		Location: cloudProps.Datacenter,
		Ram:      cloudProps.Ram,
		DiskSize: cloudProps.DiskSize,
		Cores:    cloudProps.Cores,
		Network:  networks.AsNetworkConfiguration(),
	}

	metadata := vm.InstanceMetadata{
		vm.NewSSHKeys(cv.connector.AuthorizedKeys()),
		vm.NewUserData(name, cv.connector.AgentRegistryEndpoint(),
			nil, nil),
	}
	instance, err := creator.CreateInstance(icfg, metadata)

	// Start a local forward ssh tunnel?
	if err == nil {
		err = instance.EnsureReachable(cv.connector, cv.logger)
	}

	if err != nil {
		return "", bosherr.WrapError(err, "Error launching new instance")
	}

	//TODO: add agent networks
	if err := cv.updateRegistry(agentID, instance.ID(), name, agentNetworks, env); err != nil {
		return "", err
	}
	return VMCID(instance.ID()), nil
}

func (cv CreateVM) vmName() string {

	suffix, err := cv.uuidGen.Generate()
	if err != nil {
		suffix = time.Now().Format(time.Stamp)
	}
	return fmt.Sprintf("bosh-%s", suffix)
}

func (cv CreateVM) updateRegistry(agentID string, instanceID string, vmName string,
	agentNetworks registry.NetworksSettings, env Environment) error {

	// Handle registry update failure. Delete VM or retry?
	var err error
	defer func() {
		if err != nil {
			cv.logger.Error(logTag, "Registry update failed! FIXME: handle failure")
		}
	}()
	agentOptions := cv.connector.AgentOptions()
	agentSettings := registry.NewAgentSettings(agentID, vmName, agentNetworks,
		registry.EnvSettings(env), agentOptions)

	// Update Registry with AgentSettings
	// for the agent (agent will find them as a HTTP source)
	if err = cv.registry.Update(instanceID, agentSettings); err != nil {
		return bosherr.WrapError(err, "Create VM. Error updating registry")
	}
	return nil

}
