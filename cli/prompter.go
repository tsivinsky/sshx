package cli

import (
	"os"

	ghPrompter "github.com/cli/go-gh/v2/pkg/prompter"
)

var Prompter = ghPrompter.New(os.Stdin, os.Stdout, os.Stderr)
