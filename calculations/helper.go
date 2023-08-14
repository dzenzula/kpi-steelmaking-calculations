package calculations

import (
	"main/models"
	"time"
)

func Sum(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		if q.Value != nil {
			res += *q.Value
		}
	}

	return res
}

func Avg(m []models.Query) float64 {
	var res float64

	for _, q := range m {
		if q.Value != nil {
			res += *q.Value
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
	date := time.Date(localTime.Year(), localTime.Month(), localTime.Day()-6, 0, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")

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
	for i := range dtn {
		unixTime1 := int64(*dtn[i].Value)
		unixTime2 := int64(*dtk[i].Value)

		time1 := time.Unix(unixTime1, 0)
		time2 := time.Unix(unixTime2, 0)

		minutesDifference := time2.Sub(time1).Minutes()
		differences = append(differences, minutesDifference)
	}

	totalDifferences := 0.0
	for _, diff := range differences {
		totalDifferences += diff
	}
	averageDifference := totalDifferences / float64(len(differences))

	return averageDifference
}
