package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	ghPrompter "github.com/cli/go-gh/v2/pkg/prompter"
	"github.com/tsivinsky/sshx/config"
)

var (
	serverName = flag.String("name", "", "server name")
	configDir  = "sshx"
	configFile = "config.json"
)

func main() {
	// parses cli flags
	flag.Parse()

	// sets user config dir
	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// sets config.json filepath default when running as CLI
	filepath := path.Join(confDir, configDir, configFile)

	// opens $HOME/.config/sshx/config.json for reading
	inFile, err := os.Open(filepath)
	if err != nil {
		// file doesn't exist, must be created
		inFile, err := os.Create(filepath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer inFile.Close()
	}
	defer inFile.Close()

	// opens $HOME/.config/sshx/config.json for writing
	outFile, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer outFile.Close()

	// creates pointer config.Config using the constructure and overriding defaults with CLI values
	conf, err := config.NewConfig(
		config.WithFileInput(inFile),
		config.WithFileOutput(outFile),
	)
	if err != nil {
		fmt.Println("entre 1")
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
		err = conf.List(prompter)
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
