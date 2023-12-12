package config

import (
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesnt exists: %s", configPath)
	}

	var cfg Config

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("cannot read config %s", err)
	}
	defer file.Close()

	yamlData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read config %s", err)
	}

	err = yaml.Unmarshal(yamlData, &cfg)
	if err != nil {
		log.Fatalf("error parsing YAML: %v\n", err)
	}
	return &cfg
}
