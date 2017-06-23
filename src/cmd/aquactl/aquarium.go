package main

import (
	"flag"
	"fmt"
	"os"
)

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
