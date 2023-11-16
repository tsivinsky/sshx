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

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	switch flag.Arg(0) {
	case "add":
		err = command.Add(conf)
		break

	case "connect":
		err = command.Connect(conf, *serverName)
		break

	case "list", "ls":
		err = command.List(conf)
		break

	case "remove", "rm":
		err = command.Remove(conf)
		break

	case "update":
		err = command.Update(conf)
		break

	default:
		err = command.Connect(conf, *serverName)
		break
	}

	if err != nil {
		log.Fatal(err)
	}
}
