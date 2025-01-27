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
			log.Println("New Commit", (*state)["COMMITID"])
		}},
	)

}
