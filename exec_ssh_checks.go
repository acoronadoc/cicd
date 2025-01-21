package cicd

import (
	"strings"

	"golang.org/x/crypto/ssh"
)

type Operation int

const (
	CheckRAM Operation = iota
	CheckDisk
	CheckCommandExist
)

type CheckSSHServer struct {
	Name      string
	Operation Operation
	Params    []string
}

type CheckSSHResult struct {
	Err    bool
	Result bool
	Values map[string]string
}

func CheckServerSSH(host string, config *ssh.ClientConfig, checks []CheckSSHServer) (int, map[string]CheckSSHResult) {
	cmd := checkServerCommand(checks)

	stdout, _, ecode := ExecSSH(host, config, cmd)

	return checkServerProcessResponse(stdout, ecode, checks)
}

func checkServerCommand(checks []CheckSSHServer) string {
	cmd := ""

	for i, check := range checks {
		if i != 0 {
			cmd += " && echo '------' && "
		}

		if check.Operation == CheckRAM {
			cmd += "(free || echo 'error')"
		} else if check.Operation == CheckDisk {
			cmd += "(df || echo 'error')"
		} else if check.Operation == CheckCommandExist {
			cmd += "((" + check.Params[0] + "  > /dev/null 2>/dev/null && echo '1') || echo '0')"
		}
	}

	return cmd
}

func checkServerProcessResponse(stdout string, ecode int, checks []CheckSSHServer) (int, map[string]CheckSSHResult) {
	r := map[string]CheckSSHResult{}

	if ecode != 0 {
		return ecode, nil
	}

	parts := strings.Split(stdout, "------")
	for i, part := range parts {
		lines := strings.Split(part, "\n")

		if len(lines) == 1 && lines[0] == "error" {

			r[checks[i].Name] = CheckSSHResult{
				Err:    true,
				Result: false,
			}

		} else if checks[i].Operation == CheckRAM && len(lines) > 1 {

			fields := strings.Fields(lines[1])

			r[checks[i].Name] = CheckSSHResult{
				Err:    false,
				Result: true,
				Values: map[string]string{
					"total":     fields[1],
					"used":      fields[2],
					"available": fields[6],
				},
			}

		} else if checks[i].Operation == CheckDisk {

			for _, line := range lines {
				fields := strings.Fields(line)
				if len(fields) > 5 && fields[5] == "/" {
					r[checks[i].Name] = CheckSSHResult{
						Err:    false,
						Result: true,
						Values: map[string]string{
							"total":     fields[1],
							"used":      fields[2],
							"available": fields[3],
						},
					}
				}
			}

		} else if checks[i].Operation == CheckCommandExist {
			r[checks[i].Name] = CheckSSHResult{
				Err:    false,
				Result: strings.TrimSpace(part) == "1",
			}
		}
	}

	return 0, r
}
