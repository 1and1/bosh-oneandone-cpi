package registry

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"io"
	"fmt"
	"github.com/pkg/sftp"
	"time"
	"bytes"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file + "/id_rsa")
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SSHDownload(ip, srcFile string, destination io.Writer, sshPath string) error {
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{PublicKeyFile(sshPath)},
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
