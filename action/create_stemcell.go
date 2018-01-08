package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"time"
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

// Run extracts the image URL from the properties and delegates to
// StemcellCreator for creating an image
func (cs CreateStemcell) Run(_ string, cloudProps StemcellCloudProperties) (StemcellCID, error) {

	cloudProps.ImageID = "753E3C1F859874AA74EB63B3302601F5"
	////todo: validate properites
	//if cloudProps.ImageSourceURL == "" {
	//	return "", bosherr.Error("ImageSourceURL must be specified in the stemcell manifest")
	//}
	//if cloudProps.ImageSourceURL == "" && cloudProps.ImageID == "" {
	//	return "", bosherr.Error("Image Id must be specified in the manifest")
	//}

	TIMEOUT = 30 * time.Second
	POLLING_INTERVAL = 5 * time.Second
	x := newStemcellCreator(cs.connector, cs.logger)
	stemcell, err := x.CreateStemcell("", "", "", 64, "753E3C1F859874AA74EB63B3302601F5")
	//c := newStemcellFinder(cs.connector, cs.logger)

	//stemcell, err := c.FindStemcell(cloudProps.ImageID)
	if err != nil {
		return "0", bosherr.WrapErrorf(err, "Finding stemcell with ID '%d'", cloudProps.ImageID)
	}

	return StemcellCID(stemcell), nil

}
