package cicd

import (
	"strings"
)

func GitLastCommit(url string, branch string) (string, int) {
	out, _, ecode := ExecuteCommand("git", []string{"ls-remote", "--heads", url, branch}, nil)
	return processGitLastCommit(branch, out, ecode)
}

func GitLastCommitSSH(url string, branch string, idrsa string) (string, int) {
	out, _, ecode := ExecuteCommand("git", []string{"ls-remote", "--heads", url, branch}, []string{"GIT_SSH_VARIANT=ssh", "GIT_SSH_COMMAND=ssh -i " + idrsa})
	return processGitLastCommit(branch, out, ecode)
}
func processGitLastCommit(branch string, out string, ecode int) (string, int) {
	if ecode != 0 {
		return "", ecode
	}

	lines := strings.Split(out, "\n")
	for _, l := range lines {
		fields := strings.Fields(l)

		if len(fields) == 2 && strings.HasPrefix(fields[1], "refs/heads/") {
			br := strings.ReplaceAll(fields[1], "refs/heads/", "")

			if br == branch {
				return fields[0], ecode
			}
		}
	}

	return "", ecode
}

func GitClone(url string) int {
	_, _, ecode := ExecuteCommand("git", []string{"clone", url}, nil)
	return ecode
}

func GitCloneSSH(url string, idrsa string) int {
	_, _, ecode := ExecuteCommand("git", []string{"clone", url}, []string{"GIT_SSH_VARIANT=ssh", "GIT_SSH_COMMAND=ssh -i " + idrsa})

	return ecode
}
