package command

import (
	"github.com/tsivinsky/sshx/cli"
	"github.com/tsivinsky/sshx/config"
)

func Add(conf *config.Config) error {
	name, err := cli.Prompter.Input("Server name: ", "")
	if err != nil {
		return err
	}

	user, err := cli.Prompter.Input("Server user: ", "root")
	if err != nil {
		return err
	}

	host, err := cli.Prompter.Input("Server host: ", "")
	if err != nil {
		return err
	}

	server := config.Server{
		Name: name,
		User: user,
		Host: host,
	}

	conf.Servers = append(conf.Servers, server)

	err = conf.Write()
	if err != nil {
		return err
	}

	return nil
}
