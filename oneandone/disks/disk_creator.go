package disks

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type diskCreator struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewCreator(c client.Connector, l boshlog.Logger) Creator {
	return &diskCreator{connector: c, logger: l}
}

type CreatorFactory func(client.Connector, boshlog.Logger) Creator

func (dc *diskCreator) CreateStorage(name string, sizeinGB int, dcId string) (*sdk.BlockStorage, error) {

	sizeinGB = 20
	dcId = "908DC2072407C94C8054610AD5A53B8C"
	req := sdk.BlockStorageRequest{
		Name:         name,
		Size:         &sizeinGB,
		DatacenterId: dcId,
	}
	_, res, err := dc.connector.Client().CreateBlockStorage(&req)

	if err != nil {
		return nil, fmt.Errorf("Error creating block storage. Reason: %s", err)
	}
	//res, err := dc.connector.Client().GetBlockStorage("2D62D5A790893E5D7CCBA0130A127A10")
	//if err != nil {
	//	return nil, fmt.Errorf("Error creating block storage. Reason: %s", err)
	//}

	//wait for block storage to be ready
	dc.connector.Client().WaitForState(res, "POWERED_ON", 10, 90)

	dc.logger.Debug(diskOperationsLogTag, "Created block storage %v", res)
	return res, nil
}
