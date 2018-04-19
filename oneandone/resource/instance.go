package resource

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshretry "github.com/cloudfoundry/bosh-utils/retrystrategy"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"time"
)

type Instance struct {
	vmId       string
	sshKeyPair string

	// May not always be known
	publicIPs  []string
	privateIPs []string
}

func NewInstance(vmId string, sshKeyPair string) *Instance {
	return &Instance{vmId: vmId, sshKeyPair: sshKeyPair}
}

func NewInstanceWithPrivateIPs(vmId string, privateIPs []string) *Instance {
	return &Instance{vmId: vmId, privateIPs: privateIPs}
}

func (in *Instance) ID() string {
	return in.vmId
}

func (in *Instance) SSHKeyPair() string {
	return in.sshKeyPair
}

func (in *Instance) EnsureReachable(c client.Connector, l boshlog.Logger) error {

	err := in.queryIPs(c, l)
	if err != nil {
		return err
	}
	return in.setupSSHTunnelToAgent(c, l)
}

func (in *Instance) queryIPs(c client.Connector, l boshlog.Logger) error {
	res, err := c.Client().ListServerIps(in.ID())
	if err != nil {
		l.Debug(logTag, "Error finding IPs %s", err)
		return err
	}
	var public []string
	for _, ip := range res {
		public = append(public, ip.Ip)
	}

	in.publicIPs = make([]string, len(public))
	copy(in.publicIPs, public)

	l.Debug(logTag, "Queried IPs, Private %v, Public %v", in.privateIPs, in.publicIPs)
	return nil
}

func (in *Instance) PublicIP(c client.Connector, l boshlog.Logger) (string, error) {
	ips, err := in.PublicIPs(c, l)
	if err != nil {
		return "", err
	}
	return ips[0], nil
}

func (in *Instance) PrivateIP(c client.Connector, l boshlog.Logger) (string, error) {
	ips, err := in.PrivateIPs(c, l)
	if err != nil {
		return "", err
	}
	return ips[0], nil
}

func (in *Instance) setupSSHTunnelToAgent(c client.Connector, l boshlog.Logger) (err error) {
	tunnel := c.SSHTunnelConfig()
	if tunnel.IsConfigured() {

		duration, _ := time.ParseDuration(tunnel.Duration)
		remotePort, _ := c.AgentOptions().MBusPort()
		//todo: mind this below
		remoteIP, err := in.remoteIP(c, l, true)
		if err != nil {
			return err
		}

		// Ensure SSHD is up
		//todo: mind this below
		retryable := NewSSHDCheckerRetryable("root", remoteIP, l)
		//retryable := NewSSHDCheckerRetryable(tunnel.User, remoteIP, l)
		strategy := boshretry.NewAttemptRetryStrategy(10, 20*time.Second, retryable, l)
		err = strategy.Try()

		// Then start the port forwarder
		if err == nil {
			retryable = NewSSHPortForwarderRetryable(tunnel.LocalPort, remotePort, remoteIP, tunnel.User,
				duration, l)
			strategy = boshretry.NewAttemptRetryStrategy(2, 2*time.Second, retryable, l)
			err = strategy.Try()
		}
		return err
	}
	return nil
}

func (in *Instance) havePublicIPs() bool {
	return in.publicIPs != nil && len(in.publicIPs) > 0
}

func (in *Instance) havePrivateIPs() bool {
	return in.privateIPs != nil && len(in.privateIPs) > 0
}

func (in *Instance) remoteIP(c client.Connector, l boshlog.Logger, public bool) (string, error) {
	if public {
		return in.PublicIP(c, l)
	} else {
		return in.PrivateIP(c, l)
	}
}

func (in *Instance) PublicIPs(c client.Connector, l boshlog.Logger) ([]string, error) {
	if !in.havePublicIPs() {
		if err := in.queryIPs(c, l); err != nil {
			return nil, err
		}
	}
	return in.publicIPs, nil
}

func (in *Instance) PrivateIPs(c client.Connector, l boshlog.Logger) ([]string, error) {
	if !in.havePrivateIPs() {
		if err := in.queryIPs(c, l); err != nil {
			return nil, err
		}
	}
	return in.privateIPs, nil
}
