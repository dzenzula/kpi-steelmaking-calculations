package configurations

import (
	"log"
	"main/logger"
	"main/models"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	GlobalConfig models.Config = initConfig()
)

func initConfig() models.Config {
	configFiles := []string{"configuration/config.yml", "configs/kpi-parameters.conf.yml", "config.yaml"}
	var configName string
	var config models.Config

	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			configName = configFile
			break
		}
	}

	data, err := os.ReadFile(configName)
	if err != nil {
		logger.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}
