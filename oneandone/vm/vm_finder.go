package vm

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
)

type Finder interface {
	FindInstance(instanceID string) (*resource.Instance, error)
}

type FinderFactory func(client.Connector, boshlog.Logger) Finder

type finder struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewFinder(c client.Connector, l boshlog.Logger) Finder {
	return &finder{connector: c, logger: l}
}

func (f *finder) FindInstance(instanceID string) (*resource.Instance, error) {

	f.logger.Debug(logTag, "Looking up details of VM %s", instanceID)

	r, err := f.connector.Client().GetServer(instanceID)
	if err != nil {
		f.logger.Debug(logTag, "Error GetInstance %s", err)
		return nil, err
	}
	return resource.NewInstance(r.Id), nil

}
