package disks

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

const diskPathPrefix = "/dev/sd"
const diskPathSuffix = "abcdefghijklmnopqrstuvwxyz"

type diskAttacherDetacher struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewAttacherDetacher(c client.Connector, l boshlog.Logger) AttacherDetacher {
	return &diskAttacherDetacher{connector: c, logger: l}
}

func (ad *diskAttacherDetacher) AttachInstanceToStorage(v *sdk.BlockStorage, in *resource.Instance) (string, error) {
	var devicePath string

	ad.logger.Debug(diskOperationsLogTag, "Attaching server %s to storage %s", in.ID(), v.Id)

	res, err := ad.connector.Client().AddBlockStorageServer(v.Id, in.ID())

	if err != nil {
		ad.logger.Error(diskOperationsLogTag, "Error attaching server %v", err)
		return devicePath, err
	}

	disks, err := ad.connector.Client().ListBlockStorages(1, 50, "", in.ID(), "")

	//wait for block storage to be ready
	ad.connector.Client().WaitForState(res, "ACTIVE", 10, 90)

	// Look up for the device index
	for index, attacheddisk := range disks {
		if attacheddisk.Id == v.Id {
			devicePath = fmt.Sprintf("%s%s", diskPathPrefix, string(diskPathSuffix[index]))
		}
	}
	return devicePath, nil
}

func (ad *diskAttacherDetacher) DetachInstanceFromStorage(v *sdk.BlockStorage, in *resource.Instance) error {

	res, err := ad.connector.Client().RemoveBlockStorageServer(v.Id, in.ID())
	if err != nil {
		ad.logger.Error(diskOperationsLogTag, "Error detaching volume %v", err)
		return err
	}

	//wait for block storage to be ready
	ad.connector.Client().WaitForState(res, "ACTIVE", 10, 90)

	return nil
}
