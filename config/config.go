package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Mongo  Mongo
	Server Server
}

type Mongo struct {
	Url  string
	Mode int
}

type Server struct {
	Port    string
	Domain  string
	WebRoot string
}

type ConfigLoadError struct {
	File string
	Err  error
}

func (c ConfigLoadError) Error() string {
	return fmt.Sprintf("Unable to load config file %s with error %s", c.File, c.Err.Error())
}

func (c *Config) Load(filepath string) (e error) {
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		e = ConfigLoadError{filepath, err}
		return
	}

	json.Unmarshal(file, c)
	return
}
