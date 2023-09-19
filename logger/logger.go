package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func InitLogger() {
	currentDate := time.Now()
	date := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location()).Format("2006-01-02")

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0777)
		if err != nil {
			logger.Fatal(err)
		}
	}

	file, err := os.OpenFile(fmt.Sprintf("logs/%s_kpi-steelmaking-calculations.log", date), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Fatal(err)
	}

	logger.SetOutput(file)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}
