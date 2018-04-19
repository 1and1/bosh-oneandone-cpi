package registry

// Client represents a BOSH Registry Client.
type Client interface {
	Delete(instanceID string) error
	Fetch(ipAddress string, sshPath string) (AgentSettings, error)
	UploadFile(ipAddress string, agentSettings AgentSettings, sshPath string) error
	RunCommand(ipAddress string, command []string, sshPath string) error
	UploadRootKeyPair(ipAddress string, sshPath string) error
}
