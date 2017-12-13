package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/bosh-oneandone-cpi/oneandone/client"
)

// CreateStemcell action handles the create_stemcell method invocation
type CreateStemcell struct {
	connector client.Connector
	logger    boshlog.Logger
}

// NewCreateStemcell creates a CreateStemcell instance
func NewCreateStemcell(c client.Connector, logger boshlog.Logger) CreateStemcell {
	return CreateStemcell{connector: c, logger: logger}
}

// Run extracts the image URL from the properties and delegates to
// StemcellCreator for creating an image
func (cs CreateStemcell) Run(_ string, cloudProps StemcellCloudProperties) (StemcellCID, error) {

	if cloudProps.ImageSourceURL == "" {
		return "", bosherr.Error("ImageSourceURL must be specified in the stemcell manifest")
	}
	if cloudProps.ImageID != "" {
		return "", bosherr.Error("Image Id must be specified in the manifest")
	}

	c := newStemcellCreator(cs.connector, cs.logger)
	id, err := c.CreateStemcell(cloudProps.ImageSourceURL, cloudProps.Name)
	if err != nil {
		return "", bosherr.WrapError(err, "Error creating stemcell")
	}
	return StemcellCID(id), nil

}
