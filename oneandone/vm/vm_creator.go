package vm

import (
	//"fmt"
	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	//"github.com/oneandone/oneandone-cloudserver-sdk-go"
	//"strings"
	//"time"
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
}

type Creator interface {
	CreateInstance(icfg InstanceConfiguration, md InstanceMetadata) (*resource.Instance, error)
}

type CreatorFactory func(client.Connector, boshlog.Logger) Creator

type creator struct {
	connector client.Connector
	logger    boshlog.Logger
}

func NewCreator(c client.Connector, l boshlog.Logger) Creator {
	return &creator{connector: c, logger: l}
}

func (cv *creator) CreateInstance(icfg InstanceConfiguration,
	md InstanceMetadata) (*resource.Instance, error) {

	return cv.launchInstance(icfg, md)
}
func (cv *creator) launchInstance(icfg InstanceConfiguration, md InstanceMetadata) (*resource.Instance, error) {

	////setup firewall ports
	//var firewallPolicy oneandone.FirewallPolicyRequest
	//var firewallId string
	//var firewallData *oneandone.FirewallPolicy
	//var err error
	//var flavorId string
	//var ipId string
	//
	//if icfg.ServerIp != "" {
	//	ips, err := cv.connector.Client().ListPublicIps()
	//	if err != nil {
	//		return nil, fmt.Errorf("Error fetching public IPs. Reason: %s", err)
	//	}
	//	for _, ip := range ips {
	//		if ip.IpAddress == icfg.ServerIp && ip.AssignedTo == nil {
	//			ipId = ip.Id
	//			break
	//		}
	//
	//	}
	//	if ipId == "" {
	//		return nil, fmt.Errorf("Could find a matching public ip address", icfg.ServerIp)
	//	}
	//}
	//
	//if icfg.InstanceFlavor != "" {
	//	instances, err := cv.connector.Client().ListFixedInstanceSizes()
	//	if err != nil {
	//		return nil, fmt.Errorf("Error fetching hardware flavor. Reason: %s", err)
	//	}
	//	for _, instance := range instances {
	//		if strings.ToUpper(instance.Name) == strings.ToUpper(icfg.InstanceFlavor) {
	//			flavorId = instance.Id
	//			break
	//		}
	//
	//	}
	//	if flavorId == "" {
	//		return nil, fmt.Errorf("Could find a matching instance flavor: %s , either provide a custom hardware configurations or a valid flavir (S,M,L,XL,XXL,3XL,4XL,5XL)", icfg.InstanceFlavor)
	//	}
	//
	//}
	//
	//if icfg.Network != nil && len(icfg.Network) > 0 {
	//	//check if subnet exists at a private network
	//	firewallPolicy.Name = fmt.Sprintf("Bosh fw %v", icfg.Name)
	//	var rules []oneandone.FirewallPolicyRule
	//	for _, n := range icfg.Network {
	//		for _, p := range n.OpenPorts {
	//			rules = append(rules, oneandone.FirewallPolicyRule{Protocol: "TCP", PortTo: p.PortTo, PortFrom: p.PortFrom, SourceIp: p.Source})
	//		}
	//	}
	//	firewallPolicy.Rules = rules
	//	if len(rules) > 0 {
	//		firewallId, firewallData, err = cv.connector.Client().CreateFirewallPolicy(&firewallPolicy)
	//		if err != nil {
	//			if err != nil {
	//				return nil, fmt.Errorf("Error creating a firewall policy with the open ports provided in the config file %s", err)
	//			}
	//		}
	//		cv.connector.Client().WaitForState(firewallData, "ACTIVE", 10, 90)
	//	}
	//}
	//
	////creating the server on 1&1
	//req := oneandone.ServerRequest{
	//	Name:    icfg.Name,
	//	SSHKey:  icfg.SSHKey,
	//	PowerOn: true,
	//	Hardware: oneandone.Hardware{
	//		FixedInsSizeId:    flavorId,
	//		Ram:               icfg.Ram,
	//		Vcores:            icfg.Cores,
	//		CoresPerProcessor: 1,
	//		Hdds: []oneandone.Hdd{
	//			{
	//				Size:   icfg.DiskSize,
	//				IsMain: true,
	//			},
	//		},
	//	},
	//	FirewallPolicyId: firewallId,
	//	DatacenterId:     icfg.DatacenterId,
	//	ApplianceId:      icfg.ImageId,
	//	IpId:             ipId,
	//}
	//_, res, err := cv.connector.Client().CreateServer(&req)
	//if err != nil {
	//	return nil, err
	//}
	//
	////wait on server to be ready
	////cv.connector.Client().WaitForState(res, "POWERED_OFF", 10, 90)
	//
	////cv.connector.Client().AssignServerPrivateNetwork(res.Id, "D522A56E643EED2479F2B73810DAF5F3")
	//
	////cv.connector.Client().StartServer(res.Id)
	//
	//cv.connector.Client().WaitForState(res, "POWERED_ON", 10, 90)
	//time.Sleep(120 * time.Second)
	res, _ := cv.connector.Client().GetServer("58E0D2AABD40BCC59B3BACBD4480CDDB")

	instance := resource.NewInstance(res.Id)
	return instance, nil
}
