package vm

type NetworkConfiguration struct {
	PrivateNetworkId string
	Subnet           string
	OpenPorts        []Rule
}
type Rule struct {
	PortFrom *int
	PortTo   *int
	Source   string
}
