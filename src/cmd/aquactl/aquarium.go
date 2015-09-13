package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type LightConf struct {
	Pin      uint
	Name     string
	Schedule Schedule
}

type Config struct {
	Lights []LightConf

	mu   sync.Mutex
	file string
}

func LoadConfig(p string) (*Config, error) {

	buf, err := ioutil.ReadFile(p)
	if err != nil && os.IsNotExist(err) {
		// file does not exist - create it
		f, err := os.Create(p)
		if err != nil {
			return nil, err
		}
		return &Config{file: p}, f.Close()
	}
	if err != nil {
		return nil, err
	}
	cfg := &Config{file: p}
	return cfg, yaml.Unmarshal(buf, cfg)
}

func (cfg *Config) Save() error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if cfg.file == "" {
		return nil
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cfg.file, buf, 0640)
}

func main() {

	var (
		cfg    = flag.String("f", "", "Config file to load/store schedule information in")
		config *Config
		err    error
	)

	config = &Config{
		Lights: []LightConf{
			{18, "Overhead", Schedule{}},
		},
		file: "t.yml",
	}

	config.Save()
	os.Exit(0)
	flag.Parse()
	if *cfg == "" {
		fmt.Fprintf(os.Stderr, "WARNING: running unpersisted. Won't save schedule state")
		config = &Config{}
	} else {
		config, err = LoadConfig(*cfg)
		fmt.Println(*cfg, config, err)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
			os.Exit(1)
		}
	}

	if err = Serve(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error running web server: %s\n", err)
		os.Exit(1)
	}

}
