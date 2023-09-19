package main

import (
	"database/sql"
	"fmt"
	calc "main/calculations"
	"main/database"
	"main/logger"
	"main/models"
	"reflect"
	"sync"
	"time"
)

var report = new(models.Report)

func main() {
	for {
		waitUntilMidnight()
		logger.InitLogger()
		logger.Info("Service started work")

		startTime := time.Now()
		date := calc.GetDate(0)
		msdb := database.ConnectMs()
		pgdb := database.ConnectPgData()
		pgdbReports := database.ConnectPgReports()

		calc.CacheInit(pgdb, date)

		report.Date = date

		calculations(pgdb, date)

		database.InsertPgReport(pgdbReports, *report)
		database.InsertMsReport(msdb, *report)

		msdb.Close()
		pgdb.Close()
		pgdbReports.Close()
		logger.Info("Calculations is done!")

		elapsedTime := time.Since(startTime)
		logger.Info("Run time: ", elapsedTime)
		fmt.Printf("Run time: %s\n", elapsedTime)
	}
}

func waitUntilMidnight() {
	currentTime := time.Now()
	targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 20, 0, 0, 0, currentTime.Location())

	if currentTime.After(targetTime) {
		targetTime = targetTime.Add(24 * time.Hour)
	}

	timeToWait := targetTime.Sub(currentTime)
	logger.Info("The calculation will be in", timeToWait)
	time.Sleep(timeToWait)
}

func calculations(pgdb *sql.DB, date string) {
	var wg sync.WaitGroup
	structType := reflect.TypeOf(*report)
	numFields := structType.NumField() - 1
	wg.Add(numFields)

	go func() {
		defer wg.Done()
		report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SiInCastIron = calc.GetSiInCastIron(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.CastIronTemperature = calc.GetCastIronTemperature(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.GoodCastIron = calc.GetGoodCastIron(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SContent = calc.GetSContent(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.IngotMelting = calc.IngotMeltingAvgWeight(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.O2Content = calc.O2Content(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.LimestoneFlow = calc.LimeFlow(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.DolomiteFlow = calc.DolomiteFlow(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.AluminumPreheating = calc.AluminumPreheating(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MixMelting = calc.MixMelting(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SiCC = calc.FeSiConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SiModel = calc.FeSiModelConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SiMnCC = calc.SiMnConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SiMnModel = calc.SiMnModelConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MnCC = calc.FeMnConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MnModel = calc.FeMnModelConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SlagTruncationRatio = calc.SlagTruncationRatio(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SlagSkimmingRatio = calc.SlagSkimmingRatio(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.CCMeltingCycle = calc.CCMeltingCycleMinutes(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.FePercentageInSlag = calc.FePercentageInSlag(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.SlagSamplingPercentage = calc.SlagSamplingPercentage(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.GoodCCOutput = calc.GoodCCOutput(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.GoodCCMNLZOutput = calc.GoodCCMNLZOutput(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.GoodIngotOutput = calc.GoodCCIngotOutput(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ProcessingTime = calc.ProcessingTime(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ArcTime = calc.ArcTime(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.LimestoneConsumption = calc.LimestoneConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.FluorsparConsumption = calc.FluorsparConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ArgonOxygenConsumption = calc.ArgonOxygenConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ElectricityConsumption = calc.ElectricityConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.ElectrodeConsumption = calc.ElectrodeConsumption(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.InletTemperature = calc.InletTemperature(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.InletOxidation = calc.InletOxidation(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.UPKSlagAnalysis = calc.UPKSlagAnalysis(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.CastingCycle = calc.CastingCycle(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.CastingSpeed = calc.CastingSpeed(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.CastingStopperSerial = calc.CastingStopperSerial(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ1OpenSerial = calc.MNLZOpenSerial(pgdb, date, 1)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ2OpenSerial = calc.MNLZOpenSerial(pgdb, date, 2)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ3OpenSerial = calc.MNLZOpenSerial(pgdb, date, 3)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ1Streams = calc.MNLZStreams(pgdb, date, 1)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ2Streams = calc.MNLZStreams(pgdb, date, 2)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ3Streams = calc.MNLZStreams(pgdb, date, 3)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ1RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 1)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ2RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 2)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ3RepackingDuration = calc.MNLZRepackingDuration(pgdb, date, 3)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ1MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 1)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ2MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 2)
	}()
	go func() {
		defer wg.Done()
		report.MNLZ3MeltTempDeviation = calc.MNLZMeltTempDeviation(pgdb, date, 3)
	}()
	go func() {
		defer wg.Done()
		report.GoodMNLZOutput = calc.GoodMNLZOutput(pgdb, date)
	}()
	go func() {
		defer wg.Done()
		report.MetalRetentionTime = calc.MetalRetentionTime(pgdb, date)
	}()

	wg.Wait()
}
