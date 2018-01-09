package disks

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type Creator interface {
	CreateStorage(name string, sizeinGB int, datacenterId string) (*sdk.BlockStorage, error)
}

type Terminator interface {
	DeleteStorage(volumeID string) error
}

type AttacherDetacher interface {
	AttachInstanceToStorage	(v *sdk.BlockStorage, in *resource.Instance) (string, error)
	DetachInstanceFromStorage(v *sdk.BlockStorage, in *resource.Instance) error
}

type Finder interface {
	FindStorage(storageId string) (*sdk.BlockStorage, error)
	FindAllAttachedStorages(instanceID string) ([]sdk.BlockStorage, error)
}

const diskOperationsLogTag = "OAODiskOperations"

type InstanceAttacherDetacherFactory func(*resource.Instance, client.Connector, boshlog.Logger) (AttacherDetacher, error)
type AttacherDetacherFactory func(c client.Connector, l boshlog.Logger) (AttacherDetacher, error)

func NewAttacherDetacherForInstance(in *resource.Instance, c client.Connector, l boshlog.Logger) (AttacherDetacher, error) {

	return NewAttacherDetacher(c, l), nil

}
