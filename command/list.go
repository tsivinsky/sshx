package command

import (
	"fmt"

	"github.com/tsivinsky/sshx/config"
)

func List(conf *config.Config) error {
	for _, server := range conf.Servers {
		fmt.Printf("%s: %s@%s\n", server.Name, server.User, server.Host)
	}

	return nil
}
