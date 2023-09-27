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

func main() {
	logger.Info("Service started work")
	logger.Debug("Service is in Debug mode")
	logger.InitLogger()	

	for {
		waitUntilMidnight()

		cacheData := cache.ReadCache()
		if cacheData.Date == "" {
			localTime := time.Now().Local()
			date := time.Date(localTime.Year(), localTime.Month(), localTime.Day() - 1, 19, 0, 0, 0, localTime.Location()).Format("2006-01-02 15:04:05")
			cacheData.Date = date
			cache.WriteCache(date)
		}

		missedDates := calc.GetMissingDates(cacheData.Date)

		for _, date := range missedDates {
			startTime := time.Now()
			msdb := database.ConnectMs()
			pgdb := database.ConnectPgData()
			pgdbReports := database.ConnectPgReports()

			calc.CacheInit(pgdb, date)

			report.Date = date

			calculations(pgdb, date)

			database.InsertPgReport(pgdbReports, *report)
			database.InsertMsReport(msdb, *report)
			cache.WriteCache(date)

			msdb.Close()
			pgdb.Close()
			pgdbReports.Close()
			logger.Info("Calculations is done!")

			elapsedTime := time.Since(startTime)
			logger.Info("Run time: ", elapsedTime)
			fmt.Printf("Run time: %s\n", elapsedTime)
		}
	}
}

func waitUntilMidnight() {
	currentTime := time.Now()
	targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 22, 0, 0, 0, currentTime.Location())

	if currentTime.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	timeToWait := targetTime.Sub(currentTime)
	logger.Info("Next calculation will be in", timeToWait)
	time.Sleep(timeToWait)
	logger.InitLogger()
}

func calculations(pgdb *sql.DB, date string) {
	numWorkers := 2
	tasks := []func(){
		func() {
			report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdb, date)
		},
		func() {
			report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdb, date)
		},
		func() {
			report.SiInCastIron = calc.GetSiInCastIron(pgdb, date)
		},
		func() {
			report.CastIronTemperature = calc.GetCastIronTemperature(pgdb, date)
		},
		func() {
			report.GoodCastIron = calc.GetGoodCastIron(pgdb, date)
		},
		func() {
			report.SContent = calc.GetSContent(pgdb, date)
		},
		func() {
			report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdb, date)
		},
		func() {
			report.IngotMelting = calc.IngotMeltingAvgWeight(pgdb, date)
		},
		func() {
			report.O2Content = calc.O2Content(pgdb, date)
		},
		func() {
			report.LimestoneFlow = calc.LimeFlow(pgdb, date)
		},
		func() {
			report.DolomiteFlow = calc.DolomiteFlow(pgdb, date)
		},
		func() {
			report.AluminumPreheating = calc.AluminumPreheating(pgdb, date)
		},
		func() {
			report.MixMelting = calc.MixMelting(pgdb, date)
		},
		func() {
			report.SiCC = calc.FeSiConsumption(pgdb, date)
		},
		func() {
			report.SiModel = calc.FeSiModelConsumption(pgdb, date)
		},
		func() {
			report.SiMnCC = calc.SiMnConsumption(pgdb, date)
		},
		func() {
			report.SiMnModel = calc.SiMnModelConsumption(pgdb, date)
		},
		func() {
			report.MnCC = calc.FeMnConsumption(pgdb, date)
		},
		func() {
			report.MnModel = calc.FeMnModelConsumption(pgdb, date)
		},
		func() {
			report.SlagTruncationRatio = calc.SlagTruncationRatio(pgdb, date)
		},
		func() {
			report.SlagSkimmingRatio = calc.SlagSkimmingRatio(pgdb, date)
		},
		func() {
			report.CCMeltingCycle = calc.CCMeltingCycleMinutes(pgdb, date)
		},
		func() {
			report.FePercentageInSlag = calc.FePercentageInSlag(pgdb, date)
		},
		func() {
			report.SlagSamplingPercentage = calc.SlagSamplingPercentage(pgdb, date)
		},
		func() {
			report.GoodCCOutput = calc.GoodCCOutput(pgdb, date)
		},
		func() {
			report.GoodCCMNLZOutput = calc.GoodCCMNLZOutput(pgdb, date)
		},
		func() {
			report.GoodIngotOutput = calc.GoodCCIngotOutput(pgdb, date)
		},
		func() {
			report.ProcessingTime = calc.ProcessingTime(pgdb, date)
		},
		func() {
			report.ArcTime = calc.ArcTime(pgdb, date)
		},
		func() {
			report.LimestoneConsumption = calc.LimestoneConsumption(pgdb, date)
		},
		func() {
			report.FluorsparConsumption = calc.FluorsparConsumption(pgdb, date)
		},
		func() {
			report.ArgonOxygenConsumption = calc.ArgonOxygenConsumption(pgdb, date)
		},
		func() {
			report.ElectricityConsumption = calc.ElectricityConsumption(pgdb, date)
		},
		func() {
			report.ElectrodeConsumption = calc.ElectrodeConsumption(pgdb, date)
		},
		func() {
			report.InletTemperature = calc.InletTemperature(pgdb, date)
		},
		func() {
			report.InletOxidation = calc.InletOxidation(pgdb, date)
		},
		func() {
			report.UPKSlagAnalysis = calc.UPKSlagAnalysis(pgdb, date)
		},
		func() {
			report.CastingCycle = calc.CastingCycle(pgdb, date)
		},
		func() {
			report.CastingSpeed = calc.CastingSpeed(pgdb, date)
		},
		func() {
			report.CastingStopperSerial = calc.CastingStopperSerial(pgdb, date)
		},
		func() {
			report.MNLZ1OpenSerial = calc.MNLZOpenSerial(pgdb, date, 1)
		},
		func() {
			report.MNLZ2OpenSerial = calc.MNLZOpenSerial(pgdb, date, 2)
		},
		func() {
			report.MNLZ3OpenSerial = calc.MNLZOpenSerial(pgdb, date, 3)
		},
		func() {
			report.MNLZ1Streams = calc.MNLZStreams(pgdb, date, 1)
		},
		func() {
			report.MNLZ2Streams = calc.MNLZStreams(pgdb, date, 2)
		},
		func() {
			report.MNLZ3Streams = calc.MNLZStreams(pgdb, date, 3)
		},
		func() {
			report.MNLZ1RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 1)
		},
		func() {
			report.MNLZ2RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 2)
		},
		func() {
			report.MNLZ3RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 3)
		},
		func() {
			report.MNLZ1MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 1)
		},
		func() {
			report.MNLZ2MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 2)
		},
		func() {
			report.MNLZ3MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 3)
		},
		func() {
			report.GoodMNLZOutput = calc.GoodMNLZOutput(pgdb, date)
		},
		func() {
			report.MetalRetentionTime = calc.MetalRetentionTime(pgdb, date)
		},
	}

	calc.ExecuteTasks(tasks, numWorkers)
}
