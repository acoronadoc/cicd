```
import "github.com/acoronadoc/cicd"
```

Golang library for AMAZING CI/CD pipelines

## What is CI/CD

CI/CD stands for Continuous Integration and Continuous Delivery. It's a set of practices in software development that aims to automate the process of building, testing, and deploying code changes. With CI/CD, developers can frequently integrate code changes into a shared repository, and those changes are automatically built, tested, and deployed to different environments. This approach helps to improve software quality, reduce the risk of errors, and accelerate time to market.

# Features

* Lightweight, flexible and fast.
* Create pipelines.
* Connect Git repositories(Check commits and clone) using username/password or certificates.
* Execute Local or SSH/SFTP commands.
* Check for remote commands availability.
* Check for remote servers resources(RAM and disk) availability.

# Road Map

* More remote "actions": Remote execution, CPU availability, ...
* More operations with repositories.
* Log exports.
* A User console in order to manage scripts, logs and servers.

# Installation

Simple install the package to your $GOPATH with the go tool from shell:

```
go get -u github.com/acoronadoc/cicd
```

# Samples

Get last commit from a Git repository:

```
/* From public repository */
commitId, ecode := GitLastCommit("https://github.com/acoronadoc/chatbot-sample.git", "main")

/* With private key */
commitId, ecode := GitLastCommitSSH("ssh://git@localhost:2221/srv/git/tmpgit", "main", "./keys/id_ed25599")
```

Execute remote command via SSH:

```
/* With user / pass credentials */
config := CreateUserPasswordSSHConfig("devops", "123456")

stdout, stderr, ecode := ExecSSH("192.168.15.33", config, "echo 'hola'")

/* With user / SSH key credentials */
config := CreateSSHKeysSSHConfig("devops", "./keys/id_ed25599")

stdout, stderr, ecode := ExecSSH("192.168.15.33", config, "echo 'hola'")
```

Send/receive file via SFTP:

```
/* Send file with user / pass credentials */
config := CreateUserPasswordSSHConfig("devops", "123456")

ecode := SFTPCopyLocalToRemote("192.168.15.33", config, "./test.txt", "/tmp/test.txt")

/* Receive file width user / SSH key credentials */
config := CreateSSHKeysSSHConfig("devops", "./keys/id_ed25599")

ecode := SFTPCopyRemoteToLocal("192.168.15.33", config, "/tmp/test.txt", "./test.txt")
```

Check for server resources and available commands:

```
config := CreateSSHKeysSSHConfig("devops", "./keys/id_ed25599")

ecode, r := CheckServerSSH(SSHHost, config, []CheckSSHServer{
		{Name: "RAM", Operation: CheckRAM},
		{Name: "Disk", Operation: CheckDisk},
		{Name: "ls-cmd", Operation: CheckCommandExist, Params: []string{"ls"}},
		{Name: "kubectl-cmd", Operation: CheckCommandExist, Params: []string{"kubectl"}},
	})
```