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

type Config struct {
	Servers    []Server `json:"servers"`
	configDir  string
	configFile string
}

// used to override default behavior
type option func(*Config) error

// used to initialize a new Config
func NewConfig(opts ...option) (*Config, error) {
	conf := &Config{
		configDir:  "sshx",
		configFile: "config.json"}
	for _, opt := range opts {
		err := opt(conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}

func (conf *Config) Load() error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(confDir, conf.configDir)); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(confDir, conf.configDir), 0777)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path.Join(confDir, conf.configDir, conf.configFile), os.O_RDWR|os.O_CREATE, 0644)
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

	err = json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}

func (conf *Config) Write() error {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(confDir, conf.configDir, conf.configFile), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
