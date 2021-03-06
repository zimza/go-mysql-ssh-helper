package internal

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh"
)

func getKeyFile(sshKeyPath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(sshKeyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(buf)
}

func SSHClient(sshHost *string, sshUser *string, sshKeyPath *string, sshPort *string) (sshtun *ssh.Client, err error) {
	key, err := getKeyFile(*sshKeyPath)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: *sshUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	sshtun, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", *sshHost, *sshPort), sshConfig)

	return sshtun, err
}
