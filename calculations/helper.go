package calculations

import (
	"main/models"
	"strconv"
	"time"
)

func Sum(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		if q.Value != nil {
			v, _ := strconv.ParseFloat(*q.Value, 64)
			res += v
		}
	}

	return res
}

func Avg(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		if q.Value != nil {
			v, _ := strconv.ParseFloat(*q.Value, 64)
			res += v
		}
	}
	res = SafeDivision(res, float64(len(m)))

	return res
}

func Len(m []models.Query) float64 {
	return float64(len(m))
}

func GetDate() string {
	currentTime := time.Now()
	localTime := currentTime.Local()
	date := time.Date(localTime.Year(), localTime.Month(), localTime.Day()-1, 19, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")

	return date
}

func SafeDivision(a, b float64) float64 {
	if b != 0 {
		return a / b
	}
	return 0.0
}

func AvgDiffDate(dtn []models.Query, dtk []models.Query) float64 {
	differences := []float64{}

	layout := "2006-01-02 15:04:05Z"

	for i := range dtn {
		time1, _ := time.Parse(layout, *dtn[i].Value)
		time2, _ := time.Parse(layout, *dtk[i].Value)

		minutesDifference := time2.Sub(time1).Minutes()
		differences = append(differences, minutesDifference)
	}

	if len(differences) == 0 {
		return 0.0
	}

	totalDifferences := 0.0
	for _, diff := range differences {
		totalDifferences += diff
	}

	averageDifference := totalDifferences / float64(len(differences))

	return averageDifference
}
