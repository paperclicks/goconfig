package goconfig

import (
	"io"
	"io/ioutil"
	"log"
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
	Values                                              []Value `yaml:"values"`
	AppLogFile, InfoLogFile, DebugLogFile, ErrorLogFile io.Writer
}

//New creates a new Config instance and initializes all the necessary configs
func New(filename string) (*Config, error) {

	//load config.yml to Config data structure
	config, err := LoadConfig(filename)
	if err != nil {
		return &Config{}, err
	}

	//initialize environment variables from filename if no ENVs were provided on startup
	InitEnvironment(filename)

	return &config, nil
}

func saveConfig(c Config, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

//LoadConfig loads a .yml config file into the Config data structure
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

//InitEnvironment makes sure to load all ENV variables from config.yml in case not ENV were provided
func InitEnvironment(filename string) {

	//if ENVIRONMENT env variable is not set or is different from false, load all configs from file
	//when in production all ENV variables will be set from rancher
	envValue, isEnvSet := os.LookupEnv("ENVIRONMENT")

	if !isEnvSet {
		log.Println("ENVIRONMENT is not set! Loading ENV variables from config.yml")

		err := SetValuesAsEnv(filename)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

	} else {
		log.Printf("ENVIRONMENT = %s ! Loading config ENV from system...\n", envValue)
	}

}
