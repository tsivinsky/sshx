package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Prompt(qs string) (string, error) {
	r := bufio.NewReader(os.Stdin)
	fmt.Print(qs)
	ans, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(ans), nil
}
