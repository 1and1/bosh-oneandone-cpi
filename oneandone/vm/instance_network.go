package vm

type NetworkConfiguration struct {
	Subnet    string
	OpenPorts []Rule
}
type Rule struct {
	PortFrom *int
	PortTo   *int
	Source   string
}
