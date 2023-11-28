package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/tsivinsky/sshx/cli"
)

type Server struct {
	Name string `json:"name"`
	User string `json:"user"`
	Host string `json:"server"`
}

type Config struct {
	Servers []Server `json:"servers"`
	output  io.Writer
	input   io.Reader
}

// used to override default behavior
type option func(*Config) error

// used to initialize a new Config
func NewConfig(opts ...option) (*Config, error) {
	conf := &Config{
		input:  os.Stdin,
		output: os.Stdout}
	// loops through options to override defaults
	for _, opt := range opts {
		err := opt(conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}

// option to override input reading mode when used via CLI
func WithFileInput(input io.Reader) option {
	return func(c *Config) error {
		if input == nil {
			return errors.New("nil input reader")
		}
		c.input = input
		return nil
	}
}

// option to override output writing mode when used via CLI
func WithFileOutput(output io.Writer) option {
	return func(c *Config) error {
		if output == nil {
			return errors.New("nil output writer")
		}
		c.output = output
		return nil
	}
}

// load the configs from file or sdtIn
func (conf *Config) Load() error {
	data, err := io.ReadAll(conf.input)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}

// writes the configs to file or stdOut
func (conf *Config) Write() error {
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}
	// write to file
	fmt.Fprintln(conf.output, string(data))
	return nil
}

func (conf *Config) Add() error {
	name, err := cli.Prompter.Input("Server name: ", "")
	if err != nil {
		return err
	}

	user, err := cli.Prompter.Input("Server user: ", "root")
	if err != nil {
		return err
	}

	host, err := cli.Prompter.Input("Server host: ", "")
	if err != nil {
		return err
	}

	server := Server{
		Name: name,
		User: user,
		Host: host,
	}

	conf.Servers = append(conf.Servers, server)

	err = conf.Write()
	if err != nil {
		return err
	}

	return nil
}

func (conf *Config) List() error {
	for _, server := range conf.Servers {
		fmt.Printf("%s: %s@%s\n", server.Name, server.User, server.Host)
	}

	return nil
}
