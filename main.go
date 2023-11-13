package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/tsivinsky/sshx/cli"
	"github.com/tsivinsky/sshx/config"
)

/*
1. I can add servers
2. I can list servers
3. I can remove servers
4. I can connect to server via ssh
*/

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
	}

	if err != nil {
		log.Fatal(err)
	}
}

func handleAddCommand(conf *config.Config) error {
	name, err := cli.Prompt("Server name: ")
	if err != nil {
		return err
	}

	user, err := cli.Prompt("Server user: ")
	if err != nil {
		return err
	}

	host, err := cli.Prompt("Server host: ")
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
	for _, s := range conf.Servers {
		if s.Name == name {
			server = &s
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
