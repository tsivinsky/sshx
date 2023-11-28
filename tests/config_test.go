package sshx_test

import (
	"fmt"
	"os"
	"testing"
)

type promptMock struct {
}

func (p promptMock) Select(prompt string, defaultValue string, options []string) (int, error) {
	return 0, nil
}

func (p promptMock) Input(prompt string, defaultValue string) (string, error) {
	return "", nil
}

func TestLoad(t *testing.T) {
	testFile, err := os.CreateTemp("", "*_test_config.json")
	if err != nil {
		t.Errorf("error creating tempFile for testing purposes")
	}
	defer os.Remove(testFile.Name())

	fmt.Println(testFile.Name())
}
