package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"time"

	"github.com/bosh-oneandone-cpi/oneandone/client"
)

// CreateStemcell action handles the create_stemcell method invocation
type CreateStemcell struct {
	connector client.Connector
	logger    boshlog.Logger
}

var (
	TIMEOUT          time.Duration
	POLLING_INTERVAL time.Duration
)

// NewCreateStemcell creates a CreateStemcell instance
func NewCreateStemcell(c client.Connector, logger boshlog.Logger) CreateStemcell {
	return CreateStemcell{connector: c, logger: logger}
}

func (cs CreateStemcell) Run(_ string, cloudProps StemcellCloudProperties) (stemcellId string, err error) {

	TIMEOUT = 30 * time.Second
	POLLING_INTERVAL = 5 * time.Second

	creator := newStemcellCreator(cs.connector, cs.logger)

	stemcell, err := creator.CreateStemcell(cloudProps.ImageID)
	if err != nil {
		return "0", bosherr.WrapErrorf(err, "Finding stemcell with ID '%d'", cloudProps.ImageID)
	}

	return stemcell, nil

}
