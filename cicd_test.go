package cicd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const SSHHost = "127.0.0.1:2222"
const SSHUser = "test"
const SSHPass = "password"
const SSHPrivKeyFile = "./test_env/sshconfig/id_ed25519"

const GitRepoURLHttps = "https://github.com/acoronadoc/chatbot-sample.git"
const GitRepoBranch = "main"

const GitRepoURLSSH = "ssh://git@localhost:2221/srv/git/tmpgit"
const GitSSHRSA = "test_env/sshconfig/id_ed25519"

func TestSSHConnection(t *testing.T) {

	config := CreateUserPasswordSSHConfig(SSHUser, SSHPass)
	stdout, stderr, ecode := ExecSSH(SSHHost, config, "echo 'hola'")

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error SSH", ecode)
		return
	}

	if stderr != "" {
		t.Error("Se esperaba stderr vacío y se obtuvo:", stderr)
	}

	if strings.TrimSpace(stdout) != "hola" {
		t.Error("Se esperaba 'hola' y se obtuvo:", stdout)
	}
}

func TestExecSSHPublicKeys(t *testing.T) {

	config := CreateSSHKeysSSHConfig(SSHUser, SSHPrivKeyFile)

	stdout, stderr, ecode := ExecSSH(SSHHost, config, "echo 'hola'")

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error SSH", ecode)
		return
	}

	if stderr != "" {
		t.Error("Se esperaba stderr vacío y se obtuvo:", stderr)
	}

	if strings.TrimSpace(stdout) != "hola" {
		t.Error("Se esperaba 'hola' y se obtuvo:", stdout)
	}
}

func TestCheckServerSSH(t *testing.T) {

	config := CreateUserPasswordSSHConfig(SSHUser, SSHPass)

	ecode, r := CheckServerSSH(SSHHost, config, []CheckSSHServer{
		{Name: "RAM", Operation: CheckRAM},
		{Name: "Disk", Operation: CheckDisk},
		{Name: "ls-cmd", Operation: CheckCommandExist, Params: []string{"ls"}},
		{Name: "kubectl-cmd", Operation: CheckCommandExist, Params: []string{"kubectl"}},
	})

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error SSH", ecode)
		return
	}

	if r["ls-cmd"].Result == false {
		t.Error("Se esperaba true y se obtuvo", r["ls-cmd"].Result)
	}

	if r["kubectl-cmd"].Result == true {
		t.Error("Se esperaba false y se obtuvo", r["kubectl-cmd"].Result)
	}
}

func TestCheckServerSSHPublicKeys(t *testing.T) {

	config := CreateSSHKeysSSHConfig(SSHUser, SSHPrivKeyFile)
	ecode, r := CheckServerSSH(SSHHost, config, []CheckSSHServer{
		{Name: "RAM", Operation: CheckRAM},
		{Name: "Disk", Operation: CheckDisk},
		{Name: "ls-cmd", Operation: CheckCommandExist, Params: []string{"ls"}},
		{Name: "kubectl-cmd", Operation: CheckCommandExist, Params: []string{"kubectl"}},
	})

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error SSH", ecode)
		return
	}

	if r["ls-cmd"].Result == false {
		t.Error("Se esperaba true y se obtuvo", r["ls-cmd"].Result)
	}

	if r["kubectl-cmd"].Result == true {
		t.Error("Se esperaba false y se obtuvo", r["kubectl-cmd"].Result)
	}
}

func TestGitLastCommit(t *testing.T) {

	commitId, ecode := GitLastCommit(GitRepoURLHttps, GitRepoBranch)

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error Git", ecode)
		return
	}

	if commitId == "" {
		t.Error("Se esperaba un commitId diferente a '' y se obtuvo:", commitId)
	}
}

func TestGitLastCommitSSH(t *testing.T) {

	commitId, ecode := GitLastCommitSSH(GitRepoURLHttps, GitRepoBranch, GitSSHRSA)

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error Git", ecode)
		return
	}

	if commitId == "" {
		t.Error("Se esperaba un commitId diferente a '' y se obtuvo:", commitId)
	}
}

func TestGitClone(t *testing.T) {

	ecode := GitClone(GitRepoURLHttps)

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error Git", ecode)
		return
	}

	ExecuteCommand("rm", []string{"-R", "-R", "chatbot-sample"}, nil)

}

func TestGitCloneSSH(t *testing.T) {
	pwd, _ := os.Getwd()
	ecode := GitCloneSSH(GitRepoURLSSH, filepath.Join(pwd, GitSSHRSA))

	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error Git", ecode)
		return
	}

	ExecuteCommand("rm", []string{"-R", "-R", "tmpgit"}, nil)
}

func TestCopySFTP(t *testing.T) {

	config := CreateSSHKeysSSHConfig(SSHUser, SSHPrivKeyFile)
	config2 := CreateUserPasswordSSHConfig(SSHUser, SSHPass)

	ecode := SFTPCopyLocalToRemote(SSHHost, config, filepath.Join(".", "cicd_test.go"), "/tmp/cicd_test.go")
	if ecode != 0 {
		t.Error("Se esperaba 0 y se obtuvo error al copiar LocalToRemote", ecode)
		return
	}

	ecode2 := SFTPCopyRemoteToLocal(SSHHost, config2, "/tmp/cicd_test.go", filepath.Join(".", "cicd_test.tmp"))
	if ecode2 != 0 {
		t.Error("Se esperaba 0 y se obtuvo error al copiar RemoteToLocal", ecode2)
		return
	}

	ExecSSH(SSHHost, config, "rm /tmp/cicd_test.go")
	ExecuteCommand("rm", []string{filepath.Join(".", "cicd_test.tmp")}, nil)

}
