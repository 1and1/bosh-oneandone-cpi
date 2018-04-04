package registry

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
	"os"
	//"os/user"
	"github.com/pkg/sftp"
	"io"
)

const httpClientLogTag = "RegistryHTTPClient"
const httpClientMaxAttemps = 5
const httpClientRetryDelay = 5
const settingsPath = "/var/vcap/bosh/user_data.json"

// HTTPClient represents a BOSH Registry Client.
type HTTPClient struct {
	options ClientOptions
	logger  boshlog.Logger
}

// NewHTTPClient creates a new BOSH Registry Client.
func NewHTTPClient(
	options ClientOptions,
	logger boshlog.Logger,
) HTTPClient {
	return HTTPClient{
		options: options,
		logger:  logger,
	}
}

// Delete deletes the instance settings for a given instance ID.
func (c HTTPClient) Delete(instanceID string) error {
	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.options.EndpointWithCredentials(), instanceID)
	c.logger.Debug(httpClientLogTag, "Deleting agent settings from registry endpoint '%s'", endpoint)

	request, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating DELETE request for registry endpoint '%s'", endpoint)
	}

	httpResponse, err := c.doRequest(request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting agent settings from registry endpoint '%s'", endpoint)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return bosherr.Errorf("Received status code '%d' when deleting agent settings from registry endpoint '%s'", httpResponse.StatusCode, endpoint)
	}

	c.logger.Debug(httpClientLogTag, "Deleted agent settings from registry endpoint '%s'", endpoint)
	return nil
}

// Fetch gets the agent settings for a given instance ID.
func (c HTTPClient) Fetch(username string, ipAddress string) (AgentSettings, error) {
	var agentEnv AgentSettings

	contents, err := c.Download(username, ipAddress, settingsPath)
	if err != nil {
		return AgentSettings{}, bosherr.WrapError(err, "Downloading agent env from virtual guestr")
	}

	err = json.Unmarshal(contents, &agentEnv)
	if err != nil {
		return AgentSettings{}, bosherr.WrapError(err, "Unmarshalling agent env")
	}

	c.logger.Debug(httpClientLogTag, "Fetched agent env: %#v", agentEnv)

	return agentEnv, nil
}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c HTTPClient) Update(instanceID string, agentSettings AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: '%#v", agentSettings)
	}

	endpoint := fmt.Sprintf("%s/instances/%s/settings", c.options.EndpointWithCredentials(), instanceID)
	c.logger.Debug(httpClientLogTag, "Updating registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)

	putPayload := bytes.NewReader(settingsJSON)
	request, err := http.NewRequest("PUT", endpoint, putPayload)
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating PUT request for registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)
	}

	httpResponse, err := c.doRequest(request)
	if err != nil {
		return bosherr.WrapErrorf(err, "Updating registry endpoint '%s' with agent settings: '%s'", endpoint, settingsJSON)
	}

	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		return bosherr.Errorf("Received status code '%d' when updating registry endpoint '%s' with agent settings: '%s'", httpResponse.StatusCode, endpoint, settingsJSON)
	}

	c.logger.Debug(httpClientLogTag, "Updated registry endpoint '%s' with agent settings '%s'", endpoint, settingsJSON)
	return nil
}

func (c HTTPClient) doRequest(request *http.Request) (httpResponse *http.Response, err error) {
	httpClient, err := c.httpClient()
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Creating HTTP Client")
	}

	retryDelay := time.Duration(httpClientRetryDelay) * time.Second
	for attempt := 0; attempt < httpClientMaxAttemps; attempt++ {
		httpResponse, err = httpClient.Do(request)
		if err == nil {
			return httpResponse, nil
		}
		c.logger.Debug(httpClientLogTag, "Performing registry HTTP call #%d got error '%v'", attempt, err)
		time.Sleep(retryDelay)
	}

	return nil, err
}

func (c HTTPClient) httpClient() (http.Client, error) {
	httpClient := http.Client{}

	if c.options.Protocol == "https" {
		certificates, err := tls.LoadX509KeyPair(c.options.TLS.CertFile, c.options.TLS.KeyFile)
		if err != nil {
			return httpClient, bosherr.WrapError(err, "Loading X509 Key Pair")
		}

		certPool := x509.NewCertPool()
		if c.options.TLS.CACertFile != "" {
			caCert, err := ioutil.ReadFile(c.options.TLS.CACertFile)
			if err != nil {
				return httpClient, bosherr.WrapError(err, "Loading CA certificate")
			}

			if !certPool.AppendCertsFromPEM(caCert) {
				return httpClient, bosherr.WrapError(err, "Invalid CA Certificate")
			}
		}

		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{certificates},
			InsecureSkipVerify: c.options.TLS.InsecureSkipVerify,
			RootCAs:            certPool,
		}

		httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	return httpClient, nil
}

// Update updates the agent settings for a given instance ID. If there are not already agent settings for the instance, it will create ones.
func (c HTTPClient) UploadFile(username string, ipAddress string, agentSettings AgentSettings) error {
	settingsJSON, err := json.Marshal(agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "Marshalling agent settings, contents: '%#v", agentSettings)
	}
	fileName := "1and1-agent-env.json"
	//creating the settings file locally
	writeFile(settingsJSON, fileName)

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{PublicKeyFile(username, ".ssh/id_rsa")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshAddress := fmt.Sprint(ipAddress, ":22")
	// Create a new SCP client
	client := scp.NewClient(sshAddress, config)

	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Second)
		// Connect to the remote server
		err = client.Connect()
		if err == nil {
			break
		}
	}
	if err != nil {
		return bosherr.WrapErrorf(err, "Couldn't establish a connection to the remote server'")
	}

	// Open a file
	f, _ := os.Open(fileName)

	// Close session after the file has been copied
	defer client.Session.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finaly, copy the file over
	// Usage: CopyFile(fileReader, remotePath, permission)

	client.CopyFile(f, "/var/vcap/bosh/user_data.json", "0655")

	return nil
}

func (c HTTPClient) RunCommand(username string, ipAddress string, commands []string) error {

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{PublicKeyFile(username, ".ssh/id_rsa")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	err := executeCmd(commands, ipAddress, "22", config)
	if err != nil {
		return bosherr.WrapError(err, fmt.Sprintf("Exceuting command %v failed with user %s and ip address %v", commands, username, ipAddress))
	}
	return nil
}

func writeFile(fileContent []byte, filename string) {
	err := ioutil.WriteFile(filename, fileContent, 0644)
	check(err)
}

func PublicKeyFile(username string, file string) ssh.AuthMethod {
	//usr, _ := user.Current()x
	//todo: make this work for root and non root users
	file = "/" + username + "/.ssh/id_rsa" //+ username + "/" + file
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (c HTTPClient) Download(user string, ipAddress string, sourcePath string) ([]byte, error) {
	c.logger.Debug(httpClientLogTag, "Downloading file at %s", sourcePath)
	buf := &bytes.Buffer{}
	err := SSHDownload(user, ipAddress, sourcePath, buf)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Download of %q failed", sourcePath)
	}

	c.logger.Debug(httpClientLogTag, "Downloaded %d bytes", buf.Len())

	return buf.Bytes(), nil
}

func SSHDownload(username, ip, srcFile string, destination io.Writer) error {
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{PublicKeyFile(username, ".ssh/id_rsa")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshAddress := fmt.Sprint(ip, ":22")
	client, err := ssh.Dial("tcp", sshAddress, config)
	if err != nil {
		return err
	}
	defer client.Close()

	sftp, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	f, err := sftp.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteTo(destination)
	if err != nil {
		return err
	}

	return nil
}

func executeCmd(commands []string, hostname string, port string, config *ssh.ClientConfig) error {
	conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
	session, err := conn.NewSession()
	if err != nil {
		for i := 0; i < 5; i++ {
			time.Sleep(35 * time.Second)
			// Connect to the remote server
			session, err = conn.NewSession()
			if err == nil {
				break
			}
		}
		if err != nil {
			return bosherr.WrapErrorf(err, "Couldn't establish a connection to the remote server'")
		}
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	fullcommand := ""

	for index, cmd := range commands {
		if index != (len(commands) - 1) {
			fullcommand = fullcommand + cmd + " && "
		} else {
			fullcommand = fullcommand + cmd
		}
	}
	err = session.Run(fullcommand)
	if err != nil {
		return err
	}
	return nil

}
