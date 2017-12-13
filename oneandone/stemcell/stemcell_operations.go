package stemcell

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type stemcellOperations struct {
	connector client.Connector
	logger    boshlog.Logger
}

func (so stemcellOperations) DeleteStemcell(stemcellID string) error {

	cs := so.connector.Client()
	_, err := cs.DeleteImage(stemcellID)
	return err

}

func (so stemcellOperations) CreateStemcell(sourceURI string, customImageName string) (stemcellID string, err error) {

	cs := so.connector.Client()
	req := oneandone.ImageRequest{
		Name:      customImageName,
		Source:    "iso",
		Url:       sourceURI,
		Frequency: "ONCE",
	}

	_, ok, err := cs.CreateImage(&req)

	if err != nil {
		return "", fmt.Errorf("Unable to create image from source %s. Reason: %s", sourceURI, err)
	}

	var image *oneandone.Image
	waiter := imageAvailableWaiter{
		connector: so.connector,
		logger:    so.logger,
		imageProvisionedHandler: func(i *oneandone.Image) {
			image = i
		},
	}

	if err = waiter.WaitFor(ok); err != nil {
		return "", err
	}

	return image.Id, nil
}

func (so stemcellOperations) FindStemcell(imageOCID string) (stemcellID string, err error) {

	image, err := queryImage(so.connector, imageOCID)

	if err != nil {
		return "", err
	}
	return image.Id, nil
}

func queryImage(connector client.Connector, imageOCID string) (*oneandone.Image, error) {

	image, err := connector.Client().GetImage(imageOCID)
	if err != nil {
		return nil, fmt.Errorf("Error finding image %s. Reason:%s", imageOCID, err)
	}
	return image, nil
}
