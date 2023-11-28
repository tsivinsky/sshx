package config

import "github.com/tsivinsky/sshx/cli"

func (conf *Config) Remove() error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	selectedServers, err := cli.Prompter.MultiSelect("Select servers to remove: ", []string{}, options)
	if err != nil {
		return err
	}

	newServers := []Server{}
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

	err = conf.Write()
	if err != nil {
		return err
	}

	return nil
}

func (conf *Config) Update() error {
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

	err = conf.Write()
	if err != nil {
		return err
	}

	return nil
}
