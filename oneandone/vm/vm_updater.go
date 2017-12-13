package vm

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Updater interface {
	UpdateInstanceName(instanceID string, name string) error
}

type UpdaterFactory func(client.Connector, boshlog.Logger) Updater

type updater struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewUpdater(c client.Connector, l boshlog.Logger) Updater {
	return &updater{connector: c, logger: l}
}

func (u *updater) UpdateInstanceName(instanceID string, name string) error {

	_, err := u.connector.Client().RenameServer(instanceID, name, name)
	if err != nil {
		//todo: log error
		return err
	}
	return nil
}
