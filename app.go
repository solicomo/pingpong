package main

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

type Config struct {
	Type string
	Host string
	Port string
}

type App struct {
	root    string
	appName string
	config  Config
}

func (self *App) Init(root, appName string) (err error) {

	self.root = root
	self.appName = appName

	err = self.initConfig()

	return
}

func (self *App) initConfig() (err error) {

	configFile := path.Join(self.root, self.appName+".json")

	configData, err := ioutil.ReadFile(configFile)

	if err != nil {
		return
	}

	err = json.Unmarshal(configData, &self.config)

	return
}
