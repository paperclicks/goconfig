package goconfig

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//Value represents a single value in config as Key:Value pair
type Value struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

//Config is a struct that holds all the key:value pairs present in the config file
type Config struct {
	Values []Value `yaml:"values"`
}

func saveConfig(c Config, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

//LoadConfig loads a .yml config file into the Services data structure
func LoadConfig(filename string) (Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

//SetValuesAsEnv sets all the values in the config file as ENV variables
func SetValuesAsEnv(filename string) error {

	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	for _, v := range config.Values {
		os.Setenv(v.Key, v.Value)
	}

	return nil
}
