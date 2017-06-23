package main

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfigSave(t *testing.T) {
	cfg := &Config{
		Lights: []LightConf{
			{18, "test-pin", Schedule{}},
		},
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatal(err)
	}

	cfgNext := &Config{}
	err = yaml.Unmarshal(buf, cfgNext)

	if err != nil {
		t.Fatal(err)
	}

}
