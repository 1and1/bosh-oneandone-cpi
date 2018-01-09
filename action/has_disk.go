package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/bosh-oneandone-cpi/oneandone/client"
)

// HasDisk action handles the has_disk request
type HasDisk struct {
	connector client.Connector
	logger    boshlog.Logger
}

// NewHasDisk creates a HasDisk instance
func NewHasDisk(c client.Connector, l boshlog.Logger) HasDisk {
	return HasDisk{connector: c, logger: l}
}

// Run queries OCI to determine if the given block storage exists
func (hd HasDisk) Run(diskCID DiskCID) (bool, error) {

	strg, err := newDiskFinder(hd.connector, hd.logger).FindStorage(string(diskCID))
	if err != nil {
		return false, bosherr.WrapError(err, "Error finding disk")
	}
	return strg.Id == string(diskCID), nil
}
