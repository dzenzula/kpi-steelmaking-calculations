package configurations

import (
	"log"
	"main/logger"
	"main/models"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	GlobalConfig models.Config
)

func InitConfig() {
	configFiles := []string{"configuration/config.yaml", "config.yaml", "config/kpi-parameters.conf.yml"}
	var configName string
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

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		log.Fatal(err)
	}
}
