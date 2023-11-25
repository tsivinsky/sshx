package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/tsivinsky/sshx/command"
	"github.com/tsivinsky/sshx/config"
)

var (
	serverName = flag.String("name", "", "server name")
)

func main() {
	// parses cli flags
	flag.Parse()

	// sets user config dir
	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// sets filepath default when running as CLI
	filepath := path.Join(confDir, "sshx", "config.json")

	// opens $HOME/.config/sshx/config.json for reading
	inFile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
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
		fmt.Fprintln(os.Stderr, err)
	}

	// loads configuration
	err = conf.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	switch flag.Arg(0) {
	case "add":
		err = command.Add(conf)
	case "connect":
		err = command.Connect(conf, *serverName)
	case "list", "ls":
		err = command.List(conf)
	case "remove", "rm":
		err = command.Remove(conf)
	case "update":
		err = command.Update(conf)
	default:
		err = command.Connect(conf, *serverName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

}
