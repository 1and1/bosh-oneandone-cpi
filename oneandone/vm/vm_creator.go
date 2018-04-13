package vm

import (
	//"fmt"
	"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/oneandone/oneandone-cloudserver-sdk-go"
	"strings"
	"time"
)

const logTag = "VMOperations"

type InstanceConfiguration struct {
	ImageId        string
	Name           string
	ServerIp       string
	DatacenterId   string
	SecondaryIps   []string
	Cores          int
	DiskSize       int
	Ram            float32
	SSHKey         string
	Network        Networks
	InstanceFlavor string
	SSHKeyPair     string
	EphemeralDisk  int
}

type Creator interface {
	CreateInstance(icfg InstanceConfiguration) (*resource.Instance, error)
}

type CreatorFactory func(client.Connector, boshlog.Logger) Creator

type creator struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewCreator(c client.Connector, l boshlog.Logger) Creator {
	return &creator{connector: c, logger: l}
}

func (cv *creator) CreateInstance(icfg InstanceConfiguration) (*resource.Instance, error) {

	return cv.launchInstance(icfg)
}
func (cv *creator) launchInstance(icfg InstanceConfiguration) (*resource.Instance, error) {

	//setup firewall ports
	var firewallPolicy oneandone.FirewallPolicyRequest
	var firewallId string
	var firewallData *oneandone.FirewallPolicy
	var err error
	var flavorId string
	var ipId string
	var hardwareFlavor oneandone.FixedInstanceInfo

	//check if a private network was provided
	if len(icfg.Network) == 0 || icfg.Network[0].PrivateNetworkId == "" {
		return nil, fmt.Errorf("please provide a valid private network ID")
	}
	if icfg.ServerIp != "" {
		ips, err := cv.connector.Client().ListPublicIps()
		if err != nil {
			return nil, fmt.Errorf("Error fetching public IPs. Reason: %s", err)
		}
		for _, ip := range ips {
			if ip.IpAddress == icfg.ServerIp && ip.AssignedTo == nil {
				ipId = ip.Id
				break
			}

		}
	}

	if icfg.InstanceFlavor != "" {
		instances, err := cv.connector.Client().ListFixedInstanceSizes()
		if err != nil {
			return nil, fmt.Errorf("Error fetching hardware flavor. Reason: %s", err)
		}
		for _, instance := range instances {
			if strings.ToUpper(instance.Name) == strings.ToUpper(icfg.InstanceFlavor) {
				hardwareFlavor = instance
				flavorId = instance.Id
				break
			}

		}
		if flavorId == "" {
			return nil, fmt.Errorf("Could find a matching instance flavor: %s , either provide a custom hardware configurations or a valid flavir (S,M,L,XL,XXL,3XL,4XL,5XL)", icfg.InstanceFlavor)
		}

	}

	if icfg.Network != nil && len(icfg.Network) > 0 {
		//check if subnet exists at a private network
		firewallPolicy.Name = fmt.Sprintf("Bosh fw %v", icfg.Name)
		var rules []oneandone.FirewallPolicyRule
		for _, n := range icfg.Network {
			for _, p := range n.OpenPorts {
				rules = append(rules, oneandone.FirewallPolicyRule{Protocol: "TCP", PortTo: p.PortTo, PortFrom: p.PortFrom, SourceIp: p.Source})
			}
		}
		firewallPolicy.Rules = rules
		if len(rules) > 0 {
			firewallId, firewallData, err = cv.connector.Client().CreateFirewallPolicy(&firewallPolicy)
			if err != nil {
				if err != nil {
					return nil, fmt.Errorf("Error creating a firewall policy with the open ports provided in the config file %s", err)
				}
			}
			cv.connector.Client().WaitForState(firewallData, "ACTIVE", 10, 90)
		}
	}
	if icfg.EphemeralDisk == 0 {
		icfg.EphemeralDisk = 20
	}
	hardwareFlavor.Hardware.Hdds = append(hardwareFlavor.Hardware.Hdds, oneandone.Hdd{Size: icfg.EphemeralDisk, IsMain: false})

	//creating the server on 1&1
	req := oneandone.ServerRequest{
		Name:    icfg.Name,
		SSHKey:  icfg.SSHKey,
		PowerOn: false,
		Hardware: oneandone.Hardware{
			Ram:               hardwareFlavor.Hardware.Ram,
			Vcores:            hardwareFlavor.Hardware.Vcores,
			CoresPerProcessor: hardwareFlavor.Hardware.CoresPerProcessor,
			Hdds:              hardwareFlavor.Hardware.Hdds,
		},
		FirewallPolicyId: firewallId,
		DatacenterId:     icfg.DatacenterId,
		ApplianceId:      icfg.ImageId,
		IpId:             ipId,
	}
	_, res, err := cv.connector.Client().CreateServer(&req)
	if err != nil {
		return nil, err
	}

	//wait on server to be ready
	cv.connector.Client().WaitForState(res, "POWERED_OFF", 10, 90)

	//PrivateNetwork setup
	_, err = cv.connector.Client().AssignServerPrivateNetwork(res.Id, icfg.Network[0].PrivateNetworkId)
	if err != nil {
		return nil, err
	}

	pn, err := cv.connector.Client().GetPrivateNetwork(icfg.Network[0].PrivateNetworkId)

	cv.connector.Client().WaitForState(pn, "ACTIVE", 20, 90)

	cv.connector.Client().StartServer(res.Id)

	cv.connector.Client().WaitForState(res, "POWERED_ON", 10, 90)

	// extra wait the server always needed a bit more time even after the correct state
	time.Sleep(1 * time.Minute)
	instance := resource.NewInstance(res.Id, icfg.SSHKeyPair)

	return instance, nil
}
