package command

import (
	"github.com/tsivinsky/sshx/cli"
	"github.com/tsivinsky/sshx/config"
)

func Update(conf *config.Config) error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	i, err := cli.Prompter.Select("Select server to update: ", "", options)
	if err != nil {
		return err
	}

	conf.Servers[i].Name, err = cli.Prompter.Input("Server name: ", conf.Servers[i].Name)
	if err != nil {
		return err
	}

	conf.Servers[i].User, err = cli.Prompter.Input("Server user: ", conf.Servers[i].User)
	if err != nil {
		return err
	}

	conf.Servers[i].Host, err = cli.Prompter.Input("Server host: ", conf.Servers[i].Host)
	if err != nil {
		return err
	}

	err = config.Write(conf)
	if err != nil {
		return err
	}

	return nil
}
