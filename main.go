package main

import (
	"fmt"
	c "main/configuration"
	"main/database"
	"main/logger"
	"time"
)

func main() {
	logger.InitLogger()
	c.InitConfig("configuration/config.yaml")

	currentTime := time.Now()
	localTime := currentTime.Local()
	date := time.Date(localTime.Year(), localTime.Month(), localTime.Day()-4, 0, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")

	fmt.Println(date)

	//msdev := database.ConnectMsDev()
	pgdev := database.ConnectPgDev()

	q := fmt.Sprintf(c.GlobalConfig.Querries.MeltdownsCasting, date)
	database.ExecuteQuery(pgdev, q)

	logger.Info("Done!")
}
