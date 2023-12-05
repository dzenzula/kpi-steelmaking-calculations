package calculations

import (
	c "main/configuration"
	"main/logger"
	"main/models"
	"strconv"
	"sync"
	"time"
)

var layout string = "2006-01-02 15:04:05"

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
	var count float64 = 0

	for _, q := range m {
		if q.Value != nil {
			v, _ := strconv.ParseFloat(*q.Value, 64)
			res += v
			count++
		}
	}
	res = SafeDivision(res, float64(count))

	return res
}

func Len(m []models.Query) float64 {
	if m == nil {
		return 0
	}
	res := float64(len(m))
	return res
}

func GetMissingWeeks(weekDateTracker string) []string {

	currentTime := time.Now()
	parsedWeek, err := time.Parse(layout, weekDateTracker)
	if err != nil {
		logger.Error("Error parsing the week date tracker:", err.Error())
		return nil
	}

	logger.Debug("Date from weekDateTracker:", parsedWeek.Format(layout))
	missingWeeks := []string{}

	for {
		nextWeek := parsedWeek.AddDate(0, 0, 7)

		if nextWeek.After(currentTime) {
			break
		}

		missingWeeks = append(missingWeeks, nextWeek.Format(layout))
		parsedWeek = nextWeek
	}

	if len(missingWeeks) > 0 {
		logger.Debug("Found missing weeks:", missingWeeks)
	} else {
		logger.Debug("No missing weeks found.")
	}

	return missingWeeks
}

func GetMissingMonths(monthDateTracker string) []string {
	currentTime := time.Now()
	parsedMonth, err := time.Parse(layout, monthDateTracker)
	if err != nil {
		logger.Error("Error parsing the month date tracker:", err.Error())
		return nil
	}

	logger.Debug("Date from monthDateTracker:", parsedMonth.Format(layout))
	missingMonths := []string{}

	for {
		nextMonth := parsedMonth.AddDate(0, 1, 0)

		if nextMonth.After(currentTime) {
			break
		}

		missingMonths = append(missingMonths, nextMonth.Format(layout))
		parsedMonth = nextMonth
	}

	if len(missingMonths) > 0 {
		logger.Debug("Found missing weeks:", missingMonths)
	} else {
		logger.Debug("No missing weeks found.")
	}

	return missingMonths
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
		if dtn[i].Value != nil && dtk[i].Value != nil {
			time1, _ := time.Parse(layout, *dtn[i].Value)
			time2, _ := time.Parse(layout, *dtk[i].Value)

			minutesDifference := time2.Sub(time1).Minutes()
			differences = append(differences, minutesDifference)
		} else {
			differences = append(differences, 0.0)
		}

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

	n := findMaxLength(data)
	averages := make([]*float64, n)
	counts := make([]int, n)

	updateAverages(data, averages, counts)

	normalizeAverages(averages, counts)

	return averages
}

func findMaxLength(data [][]*float64) int {
	maxLength := len(data[0])
	for _, d := range data {
		if len(d) > maxLength {
			maxLength = len(d)
		}
	}
	return maxLength
}

func updateAverages(data [][]*float64, averages []*float64, counts []int) {
	for _, row := range data {
		for i, val := range row {
			if val != nil {
				updateAverage(i, val, averages, counts)
			}
		}
	}
}

func updateAverage(index int, value *float64, averages []*float64, counts []int) {
	if averages[index] == nil {
		averages[index] = new(float64)
	}
	*averages[index] += float64(*value)
	counts[index]++
}

func normalizeAverages(averages []*float64, counts []int) {
	for i := range averages {
		if averages[i] != nil {
			*averages[i] /= float64(counts[i])
		}
	}
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

func ExecuteTasks(tasks []func(), numWorkers int) {
	numFields := len(tasks)
	taskChan := make(chan func(), numFields)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				task() // Выполняем задачу
			}
		}()
	}

	for _, task := range tasks {
		task := task
		taskChan <- func() {
			task()
		}
	}

	close(taskChan)
	wg.Wait()
}
