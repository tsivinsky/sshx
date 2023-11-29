package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func (conf *Config) Add(p Prompter) error {
	name, err := p.Input("Server name: ", "")
	if err != nil {
		return err
	}

	user, err := p.Input("Server user: ", "root")
	if err != nil {
		return err
	}

	host, err := p.Input("Server host: ", "")
	if err != nil {
		return err
	}

	server := Server{
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

func (conf *Config) List(output io.Writer) error {
	for _, server := range conf.Servers {
		fmt.Fprintf(output, "%s: %s@%s\n", server.Name, server.User, server.Host)
	}

	return nil
}

func (conf *Config) Remove(p Prompter) error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	selectedServers, err := p.MultiSelect("Select servers to remove: ", []string{}, options)
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

func (conf *Config) Update(p Prompter) error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	i, err := p.Select("Select server to update: ", "", options)
	if err != nil {
		return err
	}

	conf.Servers[i].Name, err = p.Input("Server name: ", conf.Servers[i].Name)
	if err != nil {
		return err
	}

	conf.Servers[i].User, err = p.Input("Server user: ", conf.Servers[i].User)
	if err != nil {
		return err
	}

	conf.Servers[i].Host, err = p.Input("Server host: ", conf.Servers[i].Host)
	if err != nil {
		return err
	}

	err = conf.Write()
	if err != nil {
		return err
	}

	return nil
}

func (conf *Config) Connect(p Prompter, name string) error {
	var server *Server

	if name == "" {
		options := []string{}
		for _, s := range conf.Servers {
			options = append(options, s.Name)
		}

		i, err := p.Select("Choose server: ", "", options)
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
		return errors.New("no server with this name")
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
