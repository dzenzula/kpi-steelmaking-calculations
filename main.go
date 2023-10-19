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

var report = new(models.Report)
var layout string = "2006-01-02 15:04:05"

func main() {
	logger.Info("Service started work")
	logger.Debug("Service is in Debug mode")
	logger.InitLogger()

	cacheData := cache.ReadCache()
	if cacheData.WeeklyDate == "" {
		localTime := time.Now().Local()
		nextMon := getNextMonday(time.Now(), localTime.Location())
		date := time.Date(localTime.Year(), localTime.Month(), nextMon.Day()-7, 0, 0, 0, 0, localTime.Location()).Format(layout)
		cacheData.WeeklyDate = date
		cache.WriteCache(nil, &date)
	}
	if cacheData.MonthDate == "" {
		localTime := time.Now().Local()
		nextMon := getNextFirstDayOfMonth(time.Now(), localTime.Location())
		date := time.Date(localTime.Year(), nextMon.Month()-1, nextMon.Day(), 0, 0, 0, 0, localTime.Location()).Format(layout)
		cacheData.MonthDate = date
		cache.WriteCache(&date, nil)
	}

	var missedDates []string = calc.GetMissingWeeks(cache.ReadCache().WeeklyDate)
	fmt.Println(missedDates)

	waitForMonday()
	waitForFirstDayOfMonth()

	// Block main thread so the program will not exit immediately
	select {}
}

func waitForMonday() {
	var missedDates []string = calc.GetMissingWeeks(cache.ReadCache().WeeklyDate)
	location := time.Local
	nextMonday := getNextMonday(time.Now(), location)

	mondayJob := func() {
		fmt.Printf("Running Monday's job at %v\n", time.Now())

		for _, date := range missedDates {
			parsedDate, _ := time.Parse(layout, date)
			parsedDate = parsedDate.AddDate(0, 0, -7)
			startDate := parsedDate.Format(layout)
			startTime := time.Now()
			msdb := database.ConnectMs()
			pgdb := database.ConnectPgData()
			pgdbReports := database.ConnectPgReports()

			calc.CacheInit(pgdb, startDate, date)

			report.Date = startDate
			report.Weekly = true

			calculations(pgdb, startDate, date)

			database.InsertPgReport(pgdbReports, *report)
			database.InsertMsReport(msdb, *report)
			cache.WriteCache(nil, &date)

			msdb.Close()
			pgdb.Close()
			pgdbReports.Close()
			logger.Info("Calculations is done!")

			elapsedTime := time.Since(startTime)
			logger.Info("Run time: ", elapsedTime)
			fmt.Printf("Run time: %s\n", elapsedTime)
		}

		waitForMonday()
	}
	time.AfterFunc(time.Until(nextMonday), mondayJob)
}

func waitForFirstDayOfMonth() {
	var missedDates []string = calc.GetMissingMonths(cache.ReadCache().MonthDate)
	location := time.Local
	nextFirstDayOfMonth := getNextFirstDayOfMonth(time.Now(), location)

	firstDayOfMonthJob := func() {
		fmt.Printf("Running first day of month's job at %v\n", time.Now())

		for _, date := range missedDates {
			parsedDate, _ := time.Parse(layout, date)
			parsedDate = parsedDate.AddDate(0, -1, 0)
			startDate := parsedDate.Format(layout)
			startTime := time.Now()
			msdb := database.ConnectMs()
			pgdb := database.ConnectPgData()
			pgdbReports := database.ConnectPgReports()

			calc.CacheInit(pgdb, startDate, date)

			report.Date = date
			report.Weekly = false

			calculations(pgdb, startDate, date)

			database.InsertPgReport(pgdbReports, *report)
			database.InsertMsReport(msdb, *report)
			cache.WriteCache(nil, &date)

			msdb.Close()
			pgdb.Close()
			pgdbReports.Close()
			logger.Info("Calculations is done!")

			elapsedTime := time.Since(startTime)
			logger.Info("Run time: ", elapsedTime)
			fmt.Printf("Run time: %s\n", elapsedTime)
		}

		waitForFirstDayOfMonth()
	}

	duration := time.Until(nextFirstDayOfMonth)
	time.AfterFunc(duration, firstDayOfMonthJob)
}

// Get the next monday's date
func getNextMonday(t time.Time, location *time.Location) time.Time {
	daysUntilMonday := (int(time.Monday) - int(t.Weekday()) + 7) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}

	nextMonday := time.Date(t.Year(), t.Month(), t.Day()+daysUntilMonday, 0, 0, 0, 0, location)

	logger.Info("Next iteration wil be in:", nextMonday.Format(layout))
	return nextMonday
}

// Get the next first day of the month's date
func getNextFirstDayOfMonth(t time.Time, location *time.Location) time.Time {
	nextMonth := t.Month() + 1
	if nextMonth > 12 {
		nextMonth = 1
	}

	nextFirstDayOfTheMonth := time.Date(t.Year(), nextMonth, 1, 0, 0, 0, 0, location)

	logger.Info("Next iteration wil be in:", nextFirstDayOfTheMonth.Format(layout))
	return nextFirstDayOfTheMonth
}

func calculations(pgdb *sql.DB, startDate string, endDate string) {
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
			report.CastIronTemperature = calc.GetCastIronTemperature(pgdb, startDate, endDate)
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
			report.O2Content = calc.O2Content(pgdb, startDate, endDate)
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
			report.CCMeltingCycle = calc.CCMeltingCycleMinutes(pgdb, startDate, endDate)
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
			report.ProcessingTime = calc.ProcessingTime(pgdb, startDate, endDate)
		},
		func() {
			report.ArcTime = calc.ArcTime(pgdb, startDate, endDate)
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
			report.InletTemperature = calc.InletTemperature(pgdb, startDate, endDate)
		},
		func() {
			report.InletOxidation = calc.InletOxidation(pgdb, startDate, endDate)
		},
		func() {
			report.UPKSlagAnalysis = calc.UPKSlagAnalysis(pgdb, startDate, endDate)
		},
		func() {
			report.CastingCycle = calc.CastingCycle(pgdb, startDate, endDate)
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
			report.MetalRetentionTime = calc.MetalRetentionTime(pgdb, startDate, endDate)
		},
	}

	calc.ExecuteTasks(tasks, numWorkers)
}
