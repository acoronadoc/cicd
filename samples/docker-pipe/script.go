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
			cicd.ExecuteCommand("docker", []string{"stop", "app1"}, []string{})
			cicd.ExecuteCommand("docker", []string{"rm", "app1"}, []string{})
			cicd.ExecuteCommand("docker", []string{"run", "-d", "--name", "app1", "-p", "8080:80", "nginx:1.27-alpine3.20"}, []string{})
			cicd.ExecuteCommand("docker", []string{"cp", "./repo/index.html", "app1:/usr/share/nginx/html/index.html"}, []string{})
		}},
	)

}
