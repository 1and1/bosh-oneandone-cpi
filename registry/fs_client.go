package registry

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
	"os"
	"bytes"
)

const fsClientLogTag = "RegistryFSClient"
const settingsPath = "/var/vcap/bosh/user_data.json"

type FSClient struct {
	options ClientOptions
	logger  boshlog.Logger
}

func NewFSClient(
	options ClientOptions,
	logger boshlog.Logger,
) FSClient {
	return FSClient{
		options: options,
		logger:  logger,
	}
}

// Delete deletes the instance settings for a given instance ID.
func (c FSClient) Delete(instanceID string) error {

	return nil
}

// Fetch gets the agent settings for a given instance ID.
func (c FSClient) Fetch(ipAddress string, sshPath string) (AgentSettings, error) {
	var agentEnv AgentSettings

	contents, err := c.Download(ipAddress, settingsPath, sshPath)
	if err != nil {
		return AgentSettings{}, bosherr.WrapError(err, "Downloading agent env from server")
	}

	err = json.Unmarshal(contents, &agentEnv)
	if err != nil {
		return AgentSettings{}, bosherr.WrapError(err, "Unmarshalling agent env")
	}

	c.logger.Debug(fsClientLogTag, "Fetched agent env: %#v", agentEnv)

	return agentEnv, nil
}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c FSClient) UploadFile(ipAddress string, agentSettings AgentSettings, sshPath string) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: '%s", agentSettings)
	}
	commands := []string{
		"> /var/vcap/bosh/user_data.json",
		fmt.Sprintf("echo '%s' >> /var/vcap/bosh/user_data.json", settingsJSON),
	}

	c.RunCommand(ipAddress, commands, sshPath)

	return nil
}

func (c FSClient) RunCommand(ipAddress string, commands []string, sshPath string) error {

	authMethod, err := PublicKeyFile(sshPath)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	err = executeCmd(commands, ipAddress, "22", config)
	if err != nil {
		return bosherr.WrapError(err, fmt.Sprintf("Exceuting command %v failed with user %s and ip address %v", commands, "root", ipAddress))
	}
	return nil
}

func (c FSClient) UploadRootKeyPair(ipAddress string, sshPath string) error {

	filesToCopy := [2]string{"id_rsa.pub", "id_rsa"}
	authMethod, err := PublicKeyFile(sshPath)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	for index, _ := range filesToCopy {

		sshAddress := fmt.Sprint(ipAddress, ":22")
		// Create a new SCP client
		client := scp.NewClient(sshAddress, config)

		for i := 0; i < 5; i++ {
			time.Sleep(10 * time.Second)
			// Connect to the remote server
			err = client.Connect()
			if err == nil {
				break
			}
		}
		if err != nil {
			return bosherr.WrapErrorf(err, "Couldn't establish a connection to the remote server'")
		}

		//copy private key
		filePath := "/" + filesToCopy[index]
		// Open a file
		f, _ := os.Open(sshPath + filePath)

		// Close session after the file has been copied
		defer client.Session.Close()

		// Close the file after it has been copied
		defer f.Close()

		// Finaly, copy the file over
		// Usage: CopyFile(fileReader, remotePath, permission)

		client.CopyFile(f, "/home/vcap/.ssh"+filePath, "0644")

	}

	return nil
}

func (c FSClient) Download(ipAddress string, sourcePath string, sshPath string) ([]byte, error) {
	c.logger.Debug(fsClientLogTag, "Downloading file at %s", sourcePath)
	buf := &bytes.Buffer{}
	err := SSHDownload(ipAddress, sourcePath, buf, sshPath)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Downloading %q failed", sourcePath)
	}

	c.logger.Debug(fsClientLogTag, "Downloaded %d bytes", buf.Len())

	return buf.Bytes(), nil
}
