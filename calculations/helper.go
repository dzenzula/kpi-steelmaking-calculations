package calculations

import (
	"main/models"
	"time"
)

func Sum(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		res += q.Value
	}

	return res
}

func Avg(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		res += q.Value
	}
	res = res / float64(len(m))

	return res
}

func Len(m []models.Query) float64 {
	return float64(len(m))
}

func GetDate() string {
	currentTime := time.Now()
	localTime := currentTime.Local()
	date := time.Date(localTime.Year(), localTime.Month(), localTime.Day()+1, 0, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")

	return date
}
