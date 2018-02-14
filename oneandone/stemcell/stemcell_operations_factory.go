package stemcell

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const stemCellLogTag = "OAOStemcell"

type Creator interface {
	CreateStemcell(imageId string) (stemcellId string, err error)
}
type CreatorFactory func(client.Connector, boshlog.Logger) Creator

type Destroyer interface {
	DeleteStemcell(stemcellId string) (err error)
}
type DestroyerFactory func(client.Connector, boshlog.Logger) Destroyer

type Finder interface {
	FindStemcell(imageOAOID string) (stemcellId string, err error)
}
type FinderFactory func(client.Connector, boshlog.Logger) Finder

func NewCreator(c client.Connector, l boshlog.Logger) Creator {
	return stemcellOperations{connector: c, logger: l}
}

func NewDestroyer(c client.Connector, l boshlog.Logger) Destroyer {
	return stemcellOperations{connector: c, logger: l}
}

func NewFinder(c client.Connector, l boshlog.Logger) Finder {
	return stemcellOperations{connector: c, logger: l}
}
