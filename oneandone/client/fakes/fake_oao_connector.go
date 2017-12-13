package fakes

import (
	"github.com/bosh-oneandone-cpi/registry"
	cclient "github.com/oneandone/oneandone-cloudserver-sdk-go"
)

type FakeConnector struct {
	AgentOptionsResult registry.AgentOptions
}

func (fc *FakeConnector) Connect() error {
	return nil
}

func (fc *FakeConnector) Client() *cclient.API {
	return nil
}

func (fc *FakeConnector) AuthorizedKeys() []string {
	return []string{"ssh-rsa-fake"}
}
func (fc *FakeConnector) AgentOptions() registry.AgentOptions {
	return fc.AgentOptionsResult
}

func (fc *FakeConnector) AgentRegistryEndpoint() string {
	return "fake-agent-registry-endpoint"
}
