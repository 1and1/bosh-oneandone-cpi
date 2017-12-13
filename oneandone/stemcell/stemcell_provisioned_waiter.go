package stemcell

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/oneandone/oneandone-cloudserver-sdk-go"
)

// imageAvailableWaiter waits for a newly created
// image to be fully provisioned
type imageAvailableWaiter struct {
	connector client.Connector
	logger    boshlog.Logger

	imageProvisionedHandler func(image *oneandone.Image)
}

func (w *imageAvailableWaiter) WaitFor(img *oneandone.Image) error {

	w.connector.Client().WaitForState(img, "ENABLED", 10, 90)
	w.logger.Debug(stemCellLogTag, "Done")
	return nil
}
