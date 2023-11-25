package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/tsivinsky/sshx/command"
	"github.com/tsivinsky/sshx/config"
)

var (
	serverName = flag.String("name", "", "server name")
)

func main() {
	flag.Parse()
	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	filepath := path.Join(confDir, "sshx", "config.json")
	inFile, err := os.Open(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer outFile.Close()

	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = conf.Load()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}
