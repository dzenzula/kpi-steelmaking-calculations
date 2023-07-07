package configurations

import (
	"io/ioutil"
	"log"
	"main/logger"
	"main/models"

	"gopkg.in/yaml.v2"
)

var (
	GlobalConfig models.Config
)

func InitConfig(configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatal(err)
	}

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatal(err)
	}
}
