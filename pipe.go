package cicd

import "time"

type RepoPipe struct {
	Action func(*map[string]interface{})
}

func StartPipe(bucle bool, pipe ...*RepoPipe) {
	state := map[string]interface{}{}

	for {
		for _, p := range pipe {
			p.Action(&state)
		}

		if !bucle {
			break
		}
	}
}

func PipeWaitForCommit(repoURL string, branch string, key string, sleepTime int) *RepoPipe {
	return &RepoPipe{
		Action: func(state *map[string]interface{}) {
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

func PipeWaitForCommitMulti(repoURL []string, branch []string, key string, sleepTime int) *RepoPipe {
	return &RepoPipe{
		Action: func(state *map[string]interface{}) {
			commitId := make([]string, len(repoURL))

			for i := range repoURL {
				commitId[i], _ = GitLastCommitSSH(repoURL[i], branch[i], key)
			}

			for {
				for i := range repoURL {
					ncid, ecode := GitLastCommitSSH(repoURL[i], branch[i], key)

					if ncid != commitId[i] && ecode == 0 {
						(*state)["COMMITID"] = ncid
						(*state)["REPOURL"] = repoURL[i]
						(*state)["BRANCH"] = branch[i]
						return
					}

					time.Sleep(time.Duration(sleepTime) * time.Second)
				}
			}
		},
	}
}
