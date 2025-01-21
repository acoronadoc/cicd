package cicd

import (
	"bytes"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

func CreateSSHKeysSSHConfig(username string, pubkeyfile string) *ssh.ClientConfig {
	key, err := os.ReadFile(pubkeyfile)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func CreateUserPasswordSSHConfig(username string, password string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func ExecSSH(host string, config *ssh.ClientConfig, command string) (string, string, int) {

	client := sshDial(host, config)
	if client == nil {
		return "", "", 4
	}

	session, err := client.NewSession()
	if err != nil {
		return "", err.Error(), 3
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("linux", 80, 40, modes); err != nil {
		return "", err.Error(), 2
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	if err != nil {
		return stdout.String(), stderr.String(), 1
	}

	return stdout.String(), stderr.String(), 0
}

func sshDial(host string, config *ssh.ClientConfig) *ssh.Client {
	if !strings.Contains(host, ":") {
		host += ":22"
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil
	}

	return client
}
