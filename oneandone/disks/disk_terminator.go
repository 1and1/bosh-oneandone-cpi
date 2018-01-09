package disks

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type diskTerminator struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewTerminator(c client.Connector, l boshlog.Logger) Terminator {
	return &diskTerminator{connector: c, logger: l}
}

type TerminatorFactory func(client.Connector, boshlog.Logger) Terminator

func (dt *diskTerminator) DeleteStorage(storageId string) error {
	strg, err := dt.connector.Client().GetBlockStorage(storageId)
	if err != nil {
		dt.logger.Error(diskOperationsLogTag, "Error deleting storage %v", err)
		return err
	}
	//wait for block storage to be ready
	dt.connector.Client().WaitForState(strg, "ACTIVE", 10, 90)
	_, errs := dt.connector.Client().DeleteBlockStorage(storageId)
	if errs != nil {
		dt.logger.Error(diskOperationsLogTag, "Error deleting storage %v", errs)
		return errs
	}
	dt.logger.Debug(diskOperationsLogTag, "Deleted storage %s", storageId)
	return nil
}
