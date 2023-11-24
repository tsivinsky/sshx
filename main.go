package main

import (
	"flag"
	"log"

	"github.com/tsivinsky/sshx/command"
	"github.com/tsivinsky/sshx/config"
)

var (
	serverName = flag.String("name", "", "server name")
)

func main() {
	flag.Parse()

	var err error
	conf := &config.Config{}
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
