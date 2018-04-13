package action

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/registry"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

// DetachDisk action handles the detach_disk request to detach
// a persistent disk from a vm instance
type DetachDisk struct {
	connector      client.Connector
	logger         boshlog.Logger
	registryClient registry.Client
}

// NewDetachDisk creates a DetachDisk instance
func NewDetachDisk(c client.Connector, l boshlog.Logger, r registry.Client) DetachDisk {
	return DetachDisk{connector: c, logger: l, registryClient: r}
}

// Run detaches the given disk from the the given vm. It also updates the agent registry
// after the detachment is completed. An error is thrown in case the disk or vm is not found,
// there is a failure in detachment, or if the registry can't be updated successfully.
func (dd DetachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {

	in, err := newVMFinder(dd.connector, dd.logger).FindInstance(string(vmCID))

	if err != nil {
		return nil, bosherr.WrapError(err, "Unable to find VM")
	}

	detacher, err := newAttacherDetacherForInstance(in, dd.connector, dd.logger)
	if err != nil {
		return nil, bosherr.WrapError(err, "Error creating detacher")
	}

	strg, err := newDiskFinder(dd.connector, dd.logger).FindStorage(string(diskCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Unable to find Volume")
	}

	publicIp, err := in.PublicIP(dd.connector, dd.logger)
	if err != nil {
		return "", bosherr.WrapError(err, "Error launching new instance")
	}

	if err := detacher.DetachInstanceFromStorage(strg, in); err != nil {
		return nil, bosherr.WrapError(err, "Error detaching volume")
	}

	// Read VM agent settings
	agentSettings, err := dd.registryClient.Fetch("root", publicIp)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	sshKeyPairPath := in.SSHKeyPair()
	if sshKeyPairPath == "" {
		sshKeyPairPath = sshPairKey
	}
	// Update VM agent settings
	newAgentSettings := agentSettings.DetachPersistentDisk(string(diskCID))
	if err = dd.registryClient.UploadFile(publicIp, newAgentSettings, sshKeyPairPath); err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	return nil, nil
}
