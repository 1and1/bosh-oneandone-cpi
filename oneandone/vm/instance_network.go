package vm

type NetworkConfiguration struct {
	PrivateNetworkId string
	PolicyName       string
	OpenPorts        []Rule
}
type Rule struct {
	PortFrom *int
	PortTo   *int
	Source   string
}

const DefaultBoshPolicy = "22,80,443,8443,25777,25555,8447,6868"
const DefaultBoshPolicyName = "DefaultBoshPolicy"
