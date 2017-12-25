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
	var firewallIdsToremove []string
	//TODO: cleanup any remainig resources like firewalls and load balancers
	// find any attached firewall policies
	//list server ips
	serverIps, err := t.connector.Client().ListServerIps(instanceID)
	if err != nil {
		t.logger.Error(logTag, "Ignoring Could not list instance IP's %s")
	}

	for _, ip := range serverIps {
		firewallIdsToremove = append(firewallIdsToremove, ip.Firewall.Id)
	}

	for _, id := range firewallIdsToremove {
		firewallpolicy, _ := t.connector.Client().GetFirewallPolicy(id)

		//remove the firewall policy if no other servers attached
		if len(firewallpolicy.ServerIps) <= 1 {
			_, err := t.connector.Client().DeleteFirewallPolicy(id)
			if err != nil {
				t.logger.Error(logTag, "Failed to remove firewallpolicy id=%s with error %s", id, err)
			}
		}
	}
	// Delete instance
	vm, err := t.connector.Client().DeleteServer(instanceID, false)
	if err != nil {
		t.logger.Info(logTag, "Ignoring error deleting instance %s")
	}

	return t.connector.Client().WaitUntilDeleted(vm)
}
