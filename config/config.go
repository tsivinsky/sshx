package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

var (
	configDir  = "sshx"
	configFile = "config.json"
)

type Prompter interface {
	Select(prompt string, defaultValue string, options []string) (int, error)
	Input(prompt string, defaultValue string) (string, error)
	MultiSelect(prompt string, defaultValues []string, options []string) ([]int, error)
}

type Server struct {
	Name string `json:"name"`
	User string `json:"user"`
	Host string `json:"server"`
}

type Config struct {
	Servers []Server `json:"servers"`
	File    string
}

// used to override default behavior
type option func(*Config) error

// used to initialize a new Config
func NewConfig(opts ...option) (*Config, error) {
	// sets user config dir
	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// sets config.json filepath default when running as CLI
	filepath := path.Join(confDir, configDir, configFile)
	conf := &Config{
		File: filepath}
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
func WithFile(file string) option {
	return func(c *Config) error {
		if file == "" {
			return errors.New("nil file path")
		}
		c.File = file
		return nil
	}
}

// load the configs from file or sdtIn
func (conf *Config) Load() error {
	file, err := os.OpenFile(conf.File, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
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
	file, err := os.OpenFile(conf.File, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	// solves file appending, there might be a better way to go around this
	// sshx add -> sshx remove to reproduce
	os.Truncate(conf.File, 0)
	defer file.Close()
	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	// write to file
	fmt.Fprint(file, string(data))
	return nil
}
