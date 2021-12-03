package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Load(filename string) (Config, error) {
	cfg := Config{}

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func LoadFromStr(str string) (Config, error) {
	cfg := Config{}

	if err := json.Unmarshal([]byte(str), &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func LoadDefault() (Config, error) {
	return Load(filepath.Dir(os.Args[0]) + "/config/application.json")
}

type WebConfig struct {
	Port    string `json:"port"`
	IsDebug bool   `json:"isDebug"`
}

type DbConfig struct {
	Type               string `json:"type"`
	Dsn                string `json:"dsn"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	MaxLifeTime        int    `json:"maxLifeTime"`
	MaxIdleTime        int    `json:"maxIdleTime"`
}

type Config struct {
	Web WebConfig           `json:"web"`
	Db  map[string]DbConfig `json:"db"`
}
