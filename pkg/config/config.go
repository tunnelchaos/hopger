package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type InfoType int

const (
	InfoTypeRSS InfoType = iota
	InfoTypeICal
	InfoPretalx
	InfoHubAssemblies
)

type Info struct {
	Name string
	URL  string
	Type InfoType
}

type Event struct {
	Name  string
	Infos []Info
}

type Server struct {
	GopherDir  string
	GopherPort string
	Hostname   string
}

type Config struct {
	Server Server
	Events []Event
}

func LoadConfig(filepath string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filepath, &config); err != nil {
		return Config{}, errors.New("Failed to load config file: " + err.Error())
	}
	return config, nil
}

func (c *Config) SaveConfig(filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return errors.New("Failed to create config file: " + err.Error())
	}
	defer f.Close()

	err = toml.NewEncoder(f).Encode(c)
	if err != nil {
		return errors.New("Failed to encode config file: " + err.Error())
	}
	return nil
}
