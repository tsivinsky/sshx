package command

import (
	"github.com/tsivinsky/sshx/cli"
	"github.com/tsivinsky/sshx/config"
)

func Remove(conf *config.Config) error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	selectedServers, err := cli.Prompter.MultiSelect("Select servers to remove: ", []string{}, options)
	if err != nil {
		return err
	}

	newServers := []config.Server{}
	for i, server := range conf.Servers {
		keep := true

		for _, j := range selectedServers {
			if i == j {
				keep = false
			}
		}

		if keep {
			newServers = append(newServers, server)
		}
	}
	conf.Servers = newServers

	err = config.Write(conf)
	if err != nil {
		return err
	}

	return nil
}
