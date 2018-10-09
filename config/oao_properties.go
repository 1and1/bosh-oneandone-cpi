package config

import (
	"fmt"
	"os"
)

// OAOProperties contains the properties for configuring
// BOSH CPI for 1&1 Cloud Infrastructure
type OAOProperties struct {
	// APIToken is the token used to connect to 1&1
	APIToken string `json:"token"`

	SSHTunnel SSHTunnel `json:"sshTunnel,omitempty"`
}

// AuthorizedKeys is the set of public
// ssh-rsa keys to be installed
// on the default initial account
// provisioned on a new vm
type AuthorizedKeys struct {
	Cpi  string `json:"cpi"`
	User string `json:"user, omitempty"`
}

// Validate raises an error if any of the mandatory
// properties are missing
func (b OAOProperties) Validate() error {

	if err := isAnyEmpty(map[string]string{
		"token": b.APIToken,
		//"cpiuser":     b.CpiUser,
		//"cpikeyfile":  b.CpiKeyFile,
	}); err != nil {
		return err
	}
	return nil
}

func isAnyEmpty(attributes map[string]string) error {
	for name, value := range attributes {
		if value == "" {
			return fmt.Errorf(" Property %s must not be empty", name)
		}
	}
	return nil
}

func validateFilePaths(paths []string) error {
	for _, path := range paths {
		if err := validateFilePath(path); err != nil {
			return err
		}
	}
	return nil
}

func validateFilePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("File %s doesn't exist", path)
	}
	return nil
}

func newSanitizedConfig(configFullPath string, b OAOProperties) OAOProperties {
	//dir := filepath.Dir(configFullPath)

	return OAOProperties{
		APIToken: b.APIToken,
		SSHTunnel: b.SSHTunnel,
	}
}
