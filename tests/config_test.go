package sshx_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/tsivinsky/sshx/config"
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

	content := []byte(`{
		"servers": [
		  {
			"name": "s1",
			"user": "u1",
			"server": "h1"
		  },
		  {
			"name": "s2",
			"user": "u2",
			"server": "h2"
		  },
		  {
			"name": "s3",
			"user": "u3",
			"server": "h3"
		  }
		]
	  }
	  `)

	testFile.Write(content)

	testConf, err := config.NewConfig(config.WithFile(testFile.Name()))
	if err != nil {
		t.Errorf("error creating testConf struct via NewConfig(): %v", err)
	}
	err = testConf.Load()
	if err != nil {
		t.Errorf("error calling testConf.Load(): %v", err)
	}

	want := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2", User: "u2", Host: "h2"},
		{Name: "s3", User: "u3", Host: "h3"}}
	if got := testConf.Servers; !reflect.DeepEqual(got, want) {
		t.Fatalf("errors loading config.\ngot: %v, want %v", got, want)
	}
}
