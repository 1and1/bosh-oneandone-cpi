package disks

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type diskFinder struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewFinder(c client.Connector, l boshlog.Logger) Finder {
	return &diskFinder{connector: c, logger: l}
}

type FinderFactory func(client.Connector, boshlog.Logger) Finder

func (f *diskFinder) FindStorage(storageId string) (*sdk.BlockStorage, error) {

	strg, err := f.connector.Client().GetBlockStorage(storageId)
	if err != nil {
		return nil, err
	}

	//wait for block storage to be ready
	f.connector.Client().WaitForState(strg, "POWERED_ON", 10, 90)
	return strg, nil
}

func (f *diskFinder) FindAllAttachedStorages(instanceID string) ([]sdk.BlockStorage, error) {

	strgs, err := f.connector.Client().ListBlockStorages(1, 50, "", instanceID)
	if err != nil {
		return nil, err
	}

	return strgs, nil
}
