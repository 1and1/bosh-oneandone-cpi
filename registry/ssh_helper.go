package registry

import (
	"bytes"
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"time"
)

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file + "/id_rsa")
	if err != nil {
		return nil, bosherr.WrapErrorf(err, fmt.Sprintf("Couldn't find the ssh key pair at local path"), file)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, fmt.Sprintf("Couldn't find the ssh key pair at local path"), file)
	}
	return ssh.PublicKeys(key), nil
}

func SSHDownload(ip, srcFile string, destination io.Writer, sshPath string) error {
	authMethod, err := PublicKeyFile(sshPath)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshAddress := fmt.Sprint(ip, ":22")
	for i := 0; i < 10; i++ {
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
		break
	}

	return nil
}

func executeCmd(commands []string, hostname string, port string, config *ssh.ClientConfig) error {
	var conn *ssh.Client
	var err error
	for i := 0; i < 5; i++ {
		conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
		if err != nil {
			time.Sleep(30 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		return bosherr.WrapErrorf(err, fmt.Sprintf("Couldn't establish an ssh dial to server %s on port %s with config %v'"), hostname, port, config)
	}

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
