package cicd

import (
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SFTPCopyLocalToRemote(host string, config *ssh.ClientConfig, src string, to string) int {

	client := sshDial(host, config)
	if client == nil {
		return 4
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return 3
	}
	defer sftpClient.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return 2
	}
	defer srcFile.Close()

	toFile, err := sftpClient.Create(to)
	if err != nil {
		return 1
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, srcFile)
	if err != nil {
		return 5
	}

	return 0
}

func SFTPCopyRemoteToLocal(host string, config *ssh.ClientConfig, src string, to string) int {

	client := sshDial(host, config)
	if client == nil {
		return 4
	}

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return 3
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(src)
	if err != nil {
		return 2
	}
	defer srcFile.Close()

	toFile, err := os.Create(to)
	if err != nil {
		return 1
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, srcFile)
	if err != nil {
		return 5
	}

	return 0
}
