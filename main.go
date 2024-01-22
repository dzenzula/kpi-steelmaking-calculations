package main

import (
	"database/sql"
	"fmt"
	"main/cache"
	calc "main/calculations"
	"main/database"
	"main/logger"
	"main/models"
	"time"
)

var weeklyReport = new(models.Report)
var monthlyReport = new(models.Report)
var yearlyReport = new(models.Report)
var layout string = "2006-01-02 15:04:05"

const week string = "week"
const month string = "month"
const year string = "year"

func main() {
	logger.InitLogger()
	logger.Info("Service started work")
	logger.Debug("Service is in Debug mode")

	localTime := time.Now().Local()
	cacheData := cache.ReadCache()
	if cacheData.WeeklyDate == "" {
		nextMon := getNextMonday(time.Now(), localTime.Location())
		date := time.Date(localTime.Year(), localTime.Month(), nextMon.Day()-7, 0, 0, 0, 0, localTime.Location()).Format(layout)
		cacheData.WeeklyDate = date
		cache.WriteCache(nil, &date, nil)
	}
	if cacheData.MonthDate == "" {
		nextMon := getNextFirstDayOfMonth(time.Now(), localTime.Location())
		date := time.Date(localTime.Year(), nextMon.Month()-1, nextMon.Day(), 0, 0, 0, 0, localTime.Location()).Format(layout)
		cacheData.MonthDate = date
		cache.WriteCache(&date, nil, nil)
	}
	if cacheData.YearDate == "" {
		date := time.Date(localTime.Year()-2, 1, 1, 0, 0, 0, 0, localTime.Location()).Format(layout)
		cacheData.YearDate = date
		cache.WriteCache(nil, nil, &date)
	}

	missedWeeks := calc.GetMissingWeeks(cacheData.WeeklyDate)
	if len(missedWeeks) > 0 {
		logger.Info("There is missed weeks: ", missedWeeks)
		job(week, missedWeeks)
	}
	missedMonths := calc.GetMissingMonths(cacheData.MonthDate)
	if len(missedMonths) > 0 {
		logger.Info("There is missed months: ", missedMonths)
		job(month, missedMonths)
	}
	missedYears := calc.GetMissingYears(cacheData.YearDate)
	if len(missedYears) > 0 {
		logger.Info("There is missed years: ", missedYears)
		job(year, missedYears)
	}

	CalculationService()
}

func CalculationService() {
	for {
		logger.InitLogger()
		updateYearJob()

		if time.Now().Weekday() == time.Monday {
			FirstDayOfWeek()
		}

		if time.Now().Day() == 1 {
			FirstDayOfMonth()
		}

		if time.Now().YearDay() == 1 {
			FirstDayOfYear()
		}

		duration := time.Until(waitDay())
		time.Sleep(duration)
	}
}

func FirstDayOfWeek() {
	for {
		cacheData := cache.ReadCache()
		var missedDates []string = calc.GetMissingWeeks(cacheData.WeeklyDate)
		if len(missedDates) > 0 {
			job(week, missedDates)
		}
	}
}

func FirstDayOfMonth() {
	for {
		cacheData := cache.ReadCache()
		var missedDates []string = calc.GetMissingMonths(cacheData.MonthDate)
		if len(missedDates) > 0 {
			job(month, missedDates)
		}
	}
}

func FirstDayOfYear() {
	for {
		cacheData := cache.ReadCache()
		var missedDates []string = calc.GetMissingYears(cacheData.YearDate)
		if len(missedDates) > 0 {
			job(month, missedDates)
		}
	}
}

func updateYearJob() {
	cacheData := cache.ReadCache()
	today := time.Now().UTC().Truncate(24 * time.Hour)
	parsedCache, _ := time.Parse(layout, cacheData.YearDate)
	todayString := today.Format(layout)

	logger.Info("Running year update job", todayString)

	pgdb := database.ConnectPg()

	if cacheData.YearId != 0 {
		yearlyReport.Id = cacheData.YearId
	}
	yearlyReport.Date = cacheData.YearDate
	numyear := parsedCache.Year()
	yearlyReport.ReportType = fmt.Sprintf("%d Year", numyear)
	calc.CacheInit(pgdb, cacheData.YearDate, todayString)
	calculations(pgdb, cacheData.YearDate, todayString, yearlyReport)
	yearlyReport.Id = database.UpdatePgReport(pgdb, *yearlyReport)

	pgdb.Close()
	logger.Info("Calculations is done!")
}

func job(jobType string, missedDates []string) {
	logger.Info(fmt.Sprintf("Running %s job at %v\n", func() string {
		if jobType == week {
			return "weekly"
		} else if jobType == month {
			return "monthly"
		}
		return "year"

	}(), time.Now()))

	for _, date := range missedDates {
		startTime := time.Now()
		parsedDate, _ := time.Parse(layout, date)

		// Определение даты старта в зависимости от типа задачи
		if jobType == week {
			parsedDate = parsedDate.AddDate(0, 0, -7)
		} else if jobType == month {
			parsedDate = parsedDate.AddDate(0, -1, 0)
		} else {
			parsedDate = parsedDate.AddDate(-1, 0, 0)
		}
		startDate := parsedDate.Format(layout)
		pgdb := database.ConnectPg()

		calc.CacheInit(pgdb, startDate, date)

		if jobType == week {
			weeklyReport.Date = startDate
			_, numweek := parsedDate.ISOWeek()
			weeklyReport.ReportType = fmt.Sprintf("%d Week", numweek)
			calculations(pgdb, startDate, date, weeklyReport)
			database.InsertPgReport(pgdb, *weeklyReport)
			cache.WriteCache(nil, &date, nil)
		} else if jobType == month {
			monthlyReport.Date = startDate
			monthlyReport.ReportType = parsedDate.Month().String()
			calculations(pgdb, startDate, date, monthlyReport)
			database.InsertPgReport(pgdb, *monthlyReport)
			cache.WriteCache(&date, nil, nil)
		} else {
			yearlyReport.Date = startDate
			numyear := parsedDate.Year()
			yearlyReport.ReportType = fmt.Sprintf("%d Year", numyear)
			calculations(pgdb, startDate, date, yearlyReport)
			database.InsertPgReport(pgdb, *yearlyReport)
			//database.InsertMsReport(msdb, *yearlyReport)
			cache.WriteCache(nil, nil, &date)
		}

		pgdb.Close()
		logger.Info("Calculations is done!")

		elapsedTime := time.Since(startTime)
		logger.Info("Run time: ", elapsedTime)
	}
}

// Get the next monday's date
func getNextMonday(t time.Time, location *time.Location) time.Time {
	daysUntilMonday := (int(time.Monday) - int(t.Weekday()) + 7) % 7
	nextYear := t.Year()
	nextMonth := t.Month() + 1
	if nextMonth > 12 {
		nextYear += 1
	}
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}

	nextMonday := time.Date(nextYear, t.Month(), t.Day()+daysUntilMonday, 1, 0, 0, 0, location)

	logger.Info("Next week iteration wil be in:", nextMonday.Format(layout))
	return nextMonday
}

// Get the next first day of the month's date
func getNextFirstDayOfMonth(t time.Time, location *time.Location) time.Time {
	nextMonth := t.Month() + 1
	nextYear := t.Year()
	if nextMonth > 12 {
		nextMonth = 1
		nextYear += 1
	}

	nextFirstDayOfTheMonth := time.Date(nextYear, nextMonth, 1, 1, 0, 0, 0, location)

	logger.Info("Next month iteration wil be in:", nextFirstDayOfTheMonth.Format(layout))
	return nextFirstDayOfTheMonth
}

func waitDay() time.Time {
	now := time.Now().Local().Truncate(24 * time.Hour)
	nextDay := now.Add(24 * time.Hour)

	duration := time.Until(nextDay)
	waitTime := time.Now().Add(duration)

	logger.Info("Time until the next iteration:", waitTime)

	return waitTime
}

func calculations(pgdb *sql.DB, startDate string, endDate string, report *models.Report) {
	numWorkers := 2
	tasks := []func(){
		func() {
			report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdb, startDate, endDate)
		},
		func() {
			report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdb, startDate, endDate)
		},
		func() {
			report.SiInCastIron = calc.GetSiInCastIron(pgdb, startDate, endDate)
		},
		func() {
			report.CastIronTemperature = int(calc.GetCastIronTemperature(pgdb, startDate, endDate))
		},
		func() {
			report.GoodCastIron = calc.GetGoodCastIron(pgdb, startDate, endDate)
		},
		func() {
			report.SContent = calc.GetSContent(pgdb, startDate, endDate)
		},
		func() {
			report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdb, startDate, endDate)
		},
		func() {
			report.IngotMelting = calc.IngotMeltingAvgWeight(pgdb, startDate, endDate)
		},
		func() {
			report.O2Content = int(calc.O2Content(pgdb, startDate, endDate))
		},
		func() {
			report.LimestoneFlow = calc.LimeFlow(pgdb, startDate, endDate)
		},
		func() {
			report.DolomiteFlow = calc.DolomiteFlow(pgdb, startDate, endDate)
		},
		func() {
			report.AluminumPreheating = calc.AluminumPreheating(pgdb, startDate, endDate)
		},
		func() {
			report.MixMelting = calc.MixMelting(pgdb, startDate, endDate)
		},
		func() {
			report.SiCC = calc.FeSiConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.SiModel = calc.FeSiModelConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.SiMnCC = calc.SiMnConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.SiMnModel = calc.SiMnModelConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.MnCC = calc.FeMnConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.MnModel = calc.FeMnModelConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.SlagTruncationRatio = calc.SlagTruncationRatio(pgdb, startDate, endDate)
		},
		func() {
			report.SlagSkimmingRatio = calc.SlagSkimmingRatio(pgdb, startDate, endDate)
		},
		func() {
			report.CCMeltingCycle = int(calc.CCMeltingCycleMinutes(pgdb, startDate, endDate))
		},
		func() {
			report.FePercentageInSlag = calc.FePercentageInSlag(pgdb, startDate, endDate)
		},
		func() {
			report.SlagSamplingPercentage = calc.SlagSamplingPercentage(pgdb, startDate, endDate)
		},
		func() {
			report.GoodCCOutput = calc.GoodCCOutput(pgdb, startDate, endDate)
		},
		func() {
			report.GoodCCMNLZOutput = calc.GoodCCMNLZOutput(pgdb, startDate, endDate)
		},
		func() {
			report.GoodIngotOutput = calc.GoodCCIngotOutput(pgdb, startDate, endDate)
		},
		func() {
			report.ProcessingTime = int(calc.ProcessingTime(pgdb, startDate, endDate))
		},
		func() {
			report.ArcTime = int(calc.ArcTime(pgdb, startDate, endDate))
		},
		func() {
			report.LimestoneConsumption = calc.LimestoneConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.FluorsparConsumption = calc.FluorsparConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.ArgonOxygenConsumption = calc.ArgonOxygenConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.ElectricityConsumption = calc.ElectricityConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.ElectrodeConsumption = calc.ElectrodeConsumption(pgdb, startDate, endDate)
		},
		func() {
			report.InletTemperature = int(calc.InletTemperature(pgdb, startDate, endDate))
		},
		func() {
			report.InletOxidation = int(calc.InletOxidation(pgdb, startDate, endDate))
		},
		func() {
			report.UPKSlagAnalysis = calc.UPKSlagAnalysis(pgdb, startDate, endDate)
		},
		func() {
			report.CastingCycle = int(calc.CastingCycle(pgdb, startDate, endDate))
		},
		func() {
			report.CastingSpeed = calc.CastingSpeed(pgdb, startDate, endDate)
		},
		func() {
			report.CastingStopperSerial = calc.CastingStopperSerial(pgdb, startDate, endDate)
			report.MNLZ1OpenSerial = calc.MNLZOpenSerial(pgdb, startDate, endDate, 1)
			report.MNLZ2OpenSerial = calc.MNLZOpenSerial(pgdb, startDate, endDate, 2)
			report.MNLZ3OpenSerial = calc.MNLZOpenSerial(pgdb, startDate, endDate, 3)
			report.MNLZ1MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, startDate, endDate, 1)
			report.MNLZ2MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, startDate, endDate, 2)
			report.MNLZ3MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, startDate, endDate, 3)
		},
		func() {
			report.MNLZ1Streams = calc.MNLZStreams(pgdb, startDate, endDate, 1)
		},
		func() {
			report.MNLZ2Streams = calc.MNLZStreams(pgdb, startDate, endDate, 2)
		},
		func() {
			report.MNLZ3Streams = calc.MNLZStreams(pgdb, startDate, endDate, 3)
		},
		func() {
			report.MNLZ1RepackingDuration = calc.MNLZRepackingDuration(pgdb, startDate, endDate, 1)
		},
		func() {
			report.MNLZ2RepackingDuration = calc.MNLZRepackingDuration(pgdb, startDate, endDate, 2)
		},
		func() {
			report.MNLZ3RepackingDuration = calc.MNLZRepackingDuration(pgdb, startDate, endDate, 3)
		},
		func() {
			report.GoodMNLZOutput = calc.GoodMNLZOutput(pgdb, startDate, endDate)
		},
		func() {
			report.MetalRetentionTime = int(calc.MetalRetentionTime(pgdb, startDate, endDate))
		},
	}

	calc.ExecuteTasks(tasks, numWorkers)
}
