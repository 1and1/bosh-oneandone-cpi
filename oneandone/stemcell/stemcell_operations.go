package stemcell

import (
	"fmt"

	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type stemcellOperations struct {
	connector client.Connector
	logger    boshlog.Logger
}

func (so stemcellOperations) DeleteStemcell(stemcellID string) error {

	//cs := so.connector.Client()
	//_, err := cs.DeleteImage(stemcellID)
	//return err
	return nil
}

func (so stemcellOperations) CreateStemcell(imageId string) (stemcellID string, err error) {

	image, err := queryImage(so.connector, imageId)

	if err != nil {
		return "", err
	}
	return image, nil

}

func (so stemcellOperations) FindStemcell(imageID string) (stemcellID string, err error) {

	image, err := queryImage(so.connector, imageID)

	if err != nil {
		return "", err
	}
	return image, nil
}

func queryImage(connector client.Connector, imageID string) (string, error) {
	var image *sdk.Image
	var sa *sdk.ServerAppliance
	var err error
	sa, err = connector.Client().GetServerAppliance(imageID)
	if err != nil {
		image, err = connector.Client().GetImage(imageID)
		if err != nil {
			return "", fmt.Errorf("Error finding image %s. Reason:%s", imageID, err)
		}
		return image.Id, nil
	}
	return sa.Id, nil

}
