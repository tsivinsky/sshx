package sshx_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/tsivinsky/sshx/config"
)

var content = []byte(`{
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

func TestLoad(t *testing.T) {
	t.Parallel()
	testFile, err := os.CreateTemp("", "*_test_config.json")

	if err != nil {
		t.Errorf("error creating tempFile for testing purposes")
	}
	defer os.Remove(testFile.Name())

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

func TestWrite(t *testing.T) {
	t.Parallel()
	testFile, err := os.CreateTemp("", "*_test_config.json")

	if err != nil {
		t.Errorf("error creating tempFile for testing purposes")
	}
	defer os.Remove(testFile.Name())

	testConf, err := config.NewConfig(config.WithFile(testFile.Name()))
	if err != nil {
		t.Errorf("error creating testConf struct via NewConfig(): %v", err)
	}

	servers := []config.Server{
		{Name: "s1", User: "u1", Host: "h1"},
		{Name: "s2", User: "u2", Host: "h2"},
		{Name: "s3", User: "u3", Host: "h3"},
		{Name: "s4", User: "u4", Host: "h4"}}
	testConf.Servers = servers

	err = testConf.Write()
	if err != nil {
		t.Errorf("error calling testConf.Write()")
	}
	var want = []byte(`{
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
    },
    {
      "name": "s4",
      "user": "u4",
      "server": "h4"
    }
  ]
}
`)
	got, err := io.ReadAll(testFile)
	if err != nil {
		t.Errorf("error reading testFile running testConf.Write()")
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("error testing config.Write()\ngot %v, want %v", got, want)
	}
}
