package stemcell

import (
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	sdk "github.com/oneandone/oneandone-cloudserver-sdk-go"
	//"net/http"
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

func (so stemcellOperations) CreateStemcell(sourceURI string, customImageName string, ostype string, architecture int, imageId string) (stemcellID string, err error) {

	//Todo: uncomment this when finished testing
	//cs := so.connector.Client()
	//var osid string
	//var imageSource string
	//
	//if imageId == "" {
	//	imageSource = "iso"
	//	imageOs, err := cs.ListImageOs()
	//	if err != nil {
	//		return "", fmt.Errorf("Unable to figure out os version from  ostype %s. Reason: %s", ostype, err)
	//	}
	//	for _, os := range imageOs {
	//		if os.OsVersion == ostype && *os.Architecture == architecture {
	//			osid = os.Id
	//			break
	//		}
	//
	//	}
	//} else {
	//	imageSource = "image"
	//	sourceURI = ""
	//}
	//req := oneandone.ImageRequest{
	//	Name:      customImageName,
	//	Source:    imageSource,
	//	Url:       sourceURI,
	//	Frequency: "ONCE",
	//	OsId:      osid,
	//	Type:      "os",
	//	ServerId:  imageId,
	//}
	//
	//_, image, err := cs.CreateImage(&req)
	//
	//if err != nil {
	//	return "", fmt.Errorf("Unable to create image from source %s. Reason: %s", sourceURI, err)
	//}
	//
	//waiter := imageAvailableWaiter{
	//	connector: so.connector,
	//	logger:    so.logger,
	//	imageProvisionedHandler: func(i *oneandone.Image) {
	//		image = i
	//	},
	//}
	//
	//if err = waiter.WaitFor(image); err != nil {
	//	return "", err
	//}
	//
	//if err = waiter.WaitFor(image); err != nil {
	//	return "", err
	//}
	//
	//if err = waiter.WaitFor(image); err != nil {
	//	return "", err
	//}
	//
	//return image.Id, nil
	return "753E3C1F859874AA74EB63B3302601F5", nil
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
