package cache

import (
	c "main/configuration"
	"main/logger"
	"os"

	"gopkg.in/yaml.v2"
)

type Cache struct {
	WeeklyDate string `yaml:"weekly_date"`
	MonthDate  string `yaml:"month_date"`
	YearDate   string `yaml:"year_date"`
	YearId     int    `yaml:"year_id"`
}

const errorMarshalMessage string = "Error marshaling structure: "
const errorWritingMessage string = "Error writing into file: "

func ReadCache() *Cache {
	filename := c.GlobalConfig.CachePath
	isFileExist(filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Error("Error reading the file: ", err.Error())
		return nil
	}

	var cache Cache
	err = yaml.Unmarshal(data, &cache)
	if err != nil {
		logger.Error("Error unmarshal structure: ", err.Error())
		os.WriteFile(filename, nil, 0644)
		return &cache
	}

	return &cache
}

func WriteCache(monthDate *string, weeklyDate *string, yearDate *string) {
	filename := c.GlobalConfig.CachePath
	isFileExist(filename)

	cache := ReadCache()
	if monthDate != nil {
		cache.MonthDate = *monthDate
	}
	if weeklyDate != nil {
		cache.WeeklyDate = *weeklyDate
	}
	if yearDate != nil {
		cache.YearDate = *yearDate
	}

	yamlData, err := yaml.Marshal(&cache)
	if err != nil {
		logger.Error(errorMarshalMessage, err.Error())
		return
	}

	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		logger.Error(errorWritingMessage, err.Error())
		return
	}

}

func WriteCacheId(id int) {
	filename := c.GlobalConfig.CachePath
	isFileExist(filename)
	cache := ReadCache()

	cache.YearId = id
	yamlData, err := yaml.Marshal(&cache)
	if err != nil {
		logger.Error(errorMarshalMessage, err.Error())
		return
	}
	err = os.WriteFile(filename, yamlData, 0644)
	if err != nil {
		logger.Error(errorWritingMessage, err.Error())
		return
	}
}

func isFileExist(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		config := &Cache{}

		yamlData, err := yaml.Marshal(config)
		if err != nil {
			logger.Error(errorMarshalMessage, err.Error())
			return err
		}

		err = os.WriteFile(filename, yamlData, 0644)
		if err != nil {
			logger.Error(errorWritingMessage, err.Error())
			return err
		}
	}

	return nil
}
