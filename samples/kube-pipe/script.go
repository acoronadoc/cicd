package main

import (
	"log"

	"github.com/acoronadoc/cicd"
)

const repoURL = "git@github.com:acoronadoc/test-repo.git"
const repoBranch = "main"
const repoKey = "./key"

func main() {

	cicd.StartPipe(true,

		cicd.PipeWaitForCommit(repoURL, repoBranch, repoKey, 10),

		&cicd.RepoPipe{Action: func(state *map[string]interface{}) {
			log.Println("Download Commit", (*state)["COMMITID"])

			cicd.ExecuteCommand("rm", []string{"-R", "-f", "./repo"}, []string{})
			cicd.GitCloneSSH(repoURL, repoBranch, "./repo", repoKey)
			cicd.ExecuteCommand("kubectl", []string{"delete", "configmap", "htmlapp"}, []string{})
			cicd.ExecuteCommand("kubectl", []string{"create", "configmap", "htmlapp", "--from-file=./repo/index.html"}, []string{})
			cicd.ExecuteCommand("kubectl", []string{"apply", "-f", "manifest.yaml", "--force"}, []string{})
		}},
	)

}
