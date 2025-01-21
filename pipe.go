package cicd

import "time"

type RepoPipe struct {
	action func(*map[string]interface{})
}

func StartPipe(bucle bool, pipe ...*RepoPipe) {
	state := map[string]interface{}{}

	for {
		for _, p := range pipe {
			p.action(&state)
		}

		if !bucle {
			break
		}
	}
}

func PipeWaitForCommit(repoURL string, branch string, key string, sleepTime int) *RepoPipe {
	return &RepoPipe{
		action: func(state *map[string]interface{}) {
			commitId, _ := GitLastCommitSSH(repoURL, branch, key)
			for {
				ncid, ecode := GitLastCommitSSH(repoURL, branch, key)

				if ncid != commitId && ecode == 0 {
					(*state)["COMMITID"] = ncid
					return
				}

				time.Sleep(time.Duration(sleepTime) * time.Second)
			}
		},
	}
}
