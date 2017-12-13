package vm

import (
	"github.com/bosh-oneandone-cpi/oneandone/client"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Terminator interface {
	TerminateInstance(instanceID string) error
}
type TerminatorFactory func(client.Connector, boshlog.Logger) Terminator

type terminator struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewTerminator(c client.Connector, l boshlog.Logger) Terminator {
	return &terminator{connector: c, logger: l}
}

func (t *terminator) TerminateInstance(instanceID string) error {

	t.logger.Info(logTag, "Deleting VM %s...", instanceID)

	//TODO: cleanup any remainig resources like firewalls and load balancers
	//// Find and detach any attached Vnics
	//ids, err := t.vnicAttachmentIDs(instanceID)
	//if err != nil {
	//	t.logger.Info(logTag, "Ignoring error finding attached Vnics %s", oci.CoreModelErrorMsg(err))
	//}
	//t.detachAllVnics(ids)

	// Delete instance
	vm, err := t.connector.Client().DeleteServer(instanceID, false)
	if err != nil {
		t.logger.Info(logTag, "Ignoring error deleting instance %s")
	}

	return t.connector.Client().WaitUntilDeleted(vm)
}
