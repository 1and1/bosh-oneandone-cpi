package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/bosh-oneandone-cpi/oneandone/client"
)

// GetDisks action handles the get_disks request
type GetDisks struct {
	connector client.Connector
	logger    boshlog.Logger
}

// NewGetDisks creates a GetDisks instance
func NewGetDisks(c client.Connector, l boshlog.Logger) GetDisks {
	return GetDisks{connector: c, logger: l}
}

// Run queries and returns the IDs of block storages attached to the given vm
func (gd GetDisks) Run(vmCID VMCID) ([]string, error) {

	storages, err := newDiskFinder(gd.connector, gd.logger).FindAllAttachedStorages(string(vmCID))

	if err != nil {
		return nil, bosherr.WrapError(err, "Error finding disks")
	}
	diskIds := []string{}
	for _, v := range storages {
		diskIds = append(diskIds, v.Id)
	}
	return diskIds, nil
}
