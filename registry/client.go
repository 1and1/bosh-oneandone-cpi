package registry

// Client represents a BOSH Registry Client.
type Client interface {
	Delete(instanceID string) error
	Fetch(username string, ipAddress string) (AgentSettings, error)
	Update(instanceID string, agentSettings AgentSettings) error
	UploadFile(username string, ipAddress string, agentSettings AgentSettings) error
	RunCommand(username string, ipAddress string, command []string) error
}
