package action

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"time"
)

const minVolumeSize = 20

// CreateDisk action handles the create_disk method invocation
type CreateDisk struct {
	connector client.Connector
	logger    boshlog.Logger
}

// NewCreateDisk creates a CreateDisk instance
func NewCreateDisk(c client.Connector, l boshlog.Logger) CreateDisk {
	return CreateDisk{connector: c, logger: l}
}

// Run creates a block storage  of the requested size
// and returns it's ID
func (cd CreateDisk) Run(size int, props DiskCloudProperties) (DiskCID, error) {

	creator := newDiskCreator(cd.connector, cd.logger)

	strgName := fmt.Sprintf("bosh-storage-%s", time.Now().Format(time.Stamp))
	strg, err := creator.CreateStorage(strgName, volumeSize(size), props.Datacenter)
	if err != nil {
		return "", bosherr.WrapError(err, "Error creating block storage")
	}
	return DiskCID(strg.Id), nil
}

func volumeSize(size int) int {
	//convert from mb to gb
	sizeGB := size / 1024
	s := int(sizeGB)
	if s < minVolumeSize {
		return minVolumeSize
	}
	return s
}
