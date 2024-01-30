package main

import (
	"flag"
	"fmt"
	"os"

	ghPrompter "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/tsivinsky/sshx/config"
)

var (
	serverName = flag.String("name", "", "server name")
)

var usage = `Usage: sshx <command>

Commands:
 - add
 - list, ls
 - rm
 - update`

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()

	// creates pointer config.Config using the constructure and overriding defaults with CLI values
	conf, err := config.NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	//
	prompter := ghPrompter.New(os.Stdin, os.Stdout, os.Stderr)
	// loads configuration
	err = conf.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	switch flag.Arg(0) {
	case "add":
		err = conf.Add(prompter)
	case "connect":
		err = conf.Connect(prompter, *serverName)
	case "list", "ls":
		err = conf.List(os.Stdout)
	case "remove", "rm":
		err = conf.Remove(prompter)
	case "update":
		err = conf.Update(prompter)
	default:
		err = conf.Connect(prompter, *serverName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
