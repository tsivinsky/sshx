package config

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

type Server struct {
	Name string `json:"name"`
	User string `json:"user"`
	Host string `json:"server"`
}

var (
	configDir  = "sshx"
	configFile = "config.json"
)

type Config struct {
	Servers []Server `json:"servers"`
}

func (c *Config) Load() error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(confDir, configDir)); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(confDir, configDir), 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path.Join(confDir, configDir, configFile), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if string(data) == "" {
		f.Write([]byte("{}"))
		data = []byte("{}")
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

func Write(conf *Config) error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(confDir, configDir, configFile), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
