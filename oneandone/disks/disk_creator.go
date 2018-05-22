package disks

import (
	"fmt"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"strings"
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

	var datacenterId string
	//fetch datacenter
	if dcId != "" {
		dcs, err := dc.connector.Client().ListDatacenters()
		if err != nil {
			return nil, fmt.Errorf("Error fetching data centers from API. Reason: %s", err)
		}
		for _, dc := range dcs {
			if strings.ToLower(dc.CountryCode) == strings.ToLower(dcId) {
				datacenterId = dc.Id
				break
			}
		}
	} else {
		return nil, fmt.Errorf("must provide a valid country code for datacenter (US,DE,GB,ES)")
	}
	req := sdk.BlockStorageRequest{
		Name:         name,
		Size:         &sizeinGB,
		DatacenterId: datacenterId,
	}
	_, res, err := dc.connector.Client().CreateBlockStorage(&req)

	if err != nil {
		return nil, fmt.Errorf("Error creating block storage. Reason: %s", err)
	}

	//wait for block storage to be ready
	dc.connector.Client().WaitForState(res, "POWERED_ON", 10, 90)

	dc.logger.Debug(diskOperationsLogTag, "Created block storage %v", res)
	return res, nil
}
