package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	ghPrompter "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/tsivinsky/sshx/config"
)

/*
1. I can add servers
2. I can list servers
3. I can remove servers
4. I can connect to server via ssh
*/

var prompter = ghPrompter.New(os.Stdin, os.Stdout, os.Stderr)

var (
	serverName = flag.String("name", "", "server name")
)

func main() {
	flag.Parse()

	var err error

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	command := flag.Arg(0)

	switch command {
	case "add":
		err = handleAddCommand(conf)
		break

	case "connect":
		err = handleConnectCommand(conf, *serverName)
		break

	case "list", "ls":
		err = handleListCommand(conf)
		break

	case "remove", "rm":
		err = handleRemoveCommand(conf)
		break

	default:
		err = handleConnectCommand(conf, *serverName)
		break
	}

	if err != nil {
		log.Fatal(err)
	}
}

func handleAddCommand(conf *config.Config) error {
	name, err := prompter.Input("Server name: ", "")
	if err != nil {
		return err
	}

	user, err := prompter.Input("Server user: ", "root")
	if err != nil {
		return err
	}

	host, err := prompter.Input("Server host: ", "")
	if err != nil {
		return err
	}

	server := config.Server{
		Name: name,
		User: user,
		Host: host,
	}

	conf.Servers = append(conf.Servers, server)

	err = config.Write(conf)
	if err != nil {
		return err
	}

	return nil
}

func handleConnectCommand(conf *config.Config, name string) error {
	var server *config.Server

	if name == "" {
		options := []string{}
		for _, s := range conf.Servers {
			options = append(options, s.Name)
		}

		i, err := prompter.Select("Choose server: ", "", options)
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

func handleListCommand(conf *config.Config) error {
	for _, server := range conf.Servers {
		fmt.Printf("%s: %s@%s\n", server.Name, server.User, server.Host)
	}

	return nil
}

func handleRemoveCommand(conf *config.Config) error {
	options := []string{}
	for _, server := range conf.Servers {
		options = append(options, server.Name)
	}

	selectedServers, err := prompter.MultiSelect("Select servers to remove: ", []string{}, options)
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
