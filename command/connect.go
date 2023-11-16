package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/tsivinsky/sshx/cli"
	"github.com/tsivinsky/sshx/config"
)

func Connect(conf *config.Config, name string) error {
	var server *config.Server

	if name == "" {
		options := []string{}
		for _, s := range conf.Servers {
			options = append(options, s.Name)
		}

		i, err := cli.Prompter.Select("Choose server: ", "", options)
		if err != nil {
			return err
		}

		server = &conf.Servers[i]
	}

	if server == nil {
		for _, s := range conf.Servers {
			if s.Name == name {
				server = &s
			}
		}
	}

	if server == nil {
		return errors.New("No server with this name")
	}

	serverHost := fmt.Sprintf("%s@%s", server.User, server.Host)

	cmd := exec.Command("ssh", serverHost)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
