package calculations

import (
	c "main/configuration"
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
	if m == nil {
		return 0
	}
	res := float64(len(m))
	return res
}

func GetDate(i int) string {
	currentTime := time.Now()
	localTime := currentTime.Local()
	date := time.Date(localTime.Year(), localTime.Month(), localTime.Day()+i, 19, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")

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

func CalculateAverages(data [][]*float64) []*float64 {
	if len(data) == 0 || len(data[0]) == 0 {
		return nil
	}

	n := len(data[0])
	averages := make([]*float64, n)
	counts := make([]int, n)

	for _, row := range data {
		for i, val := range row {
			if val != nil {
				if averages[i] == nil {
					averages[i] = new(float64)
				}
				*averages[i] += float64(*val)
				counts[i]++
			}
		}
	}

	for i := 0; i < n; i++ {
		if averages[i] != nil {
			*averages[i] /= float64(counts[i])
		}
	}

	return averages
}

func ParseFloatValues(queries []models.Query) []*float64 {
	values := make([]*float64, len(queries))
	for i, query := range queries {
		if query.Value != nil {
			value, err := strconv.ParseFloat(*query.Value, 64)
			if err == nil {
				values[i] = &value
			} else {
				values[i] = nil
			}
		}
	}
	return values
}

func FindSteelGrade(steelType string) (int, int) {
	for _, mark := range c.GlobalConfig.SteelMarks {
		if mark.SteelType == steelType {
			return mark.Min, mark.Max
		}
	}
	return 0, 0
}