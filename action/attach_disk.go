package action

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/registry"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const diskActionsLogTag = "diskActions"

// AttachDisk action handles the attach_disk request to attach
// a persistent disk to a vm instance
type AttachDisk struct {
	connector      client.Connector
	logger         boshlog.Logger
	registryClient registry.Client
}

// NewAttachDisk creates an AttachDisk instance
func NewAttachDisk(c client.Connector, l boshlog.Logger, r registry.Client) AttachDisk {
	return AttachDisk{connector: c, logger: l, registryClient: r}

}

// Run implements the Action handler
func (ad AttachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {

	var devicePath string
	in, err := newVMFinder(ad.connector, ad.logger).FindInstance(string(vmCID))

	if err != nil {
		return nil, bosherr.WrapError(err, "Unable to find VM")
	}

	attacher, err := newAttacherDetacherForInstance(in, ad.connector, ad.logger)

	if err != nil {
		return nil, bosherr.WrapError(err, "Error creating attacher")
	}

	vol, err := newDiskFinder(ad.connector, ad.logger).FindStorage(string(diskCID))
	if err != nil {
		return nil, bosherr.WrapError(err, "Unable to find storage")
	}

	devicePath, err = attacher.AttachInstanceToStorage(vol, in)
	if err != nil {
		if err == nil {
			err = fmt.Errorf("storage not attached %v", *vol)
		}
		return nil, bosherr.WrapError(err, "Error attaching storage")
	}

	publicIp, err := in.PublicIP(ad.connector, ad.logger)
	if err != nil {
		return "", bosherr.WrapError(err, "Error launching new instance")
	}
	//todo: sdb with devicepath value
	//defining partionion label
	bsCommands := []string{"parted -s /dev/sdb mklabel gpt"}

	ad.registryClient.RunCommand("root", publicIp, bsCommands)

	// Read VM agent settings
	agentSettings, err := ad.registryClient.Fetch("root", publicIp)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.AttachPersistentDisk(string(diskCID), devicePath)
	if err = ad.registryClient.UploadFile("root", publicIp, newAgentSettings); err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}
	return nil, nil
}
