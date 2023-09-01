package main

import (
	calc "main/calculations"
	"main/database"
	"main/logger"
	"main/models"
	"time"
)

var report = new(models.Report)

func main() {
	logger.InitLogger()

	for {
		waitUntilMidnight()

		date := calc.GetDate()
		msdb := database.ConnectMs()
		pgdb := database.ConnectPg()

		calc.CacheInit(pgdb, date)

		report.Date = date
		report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdb, date)
		report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdb, date)
		report.SiInCastIron = calc.GetSiInCastIron(pgdb, date)
		report.CastIronTemperature = calc.GetCastIronTemperature(pgdb, date)
		report.GoodCastIron = calc.GetGoodCastIron(pgdb, date)
		report.SContent = calc.GetSContent(pgdb, date)
		report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdb, date)
		report.IngotMelting = calc.IngotMeltingAvgWeight(pgdb, date)
		report.O2Content = calc.O2Content(pgdb, date)
		report.LimestoneFlow = calc.LimeFlow(pgdb, date)
		report.DolomiteFlow = calc.DolomiteFlow(pgdb, date)
		report.AluminumPreheating = calc.AluminumPreheating(pgdb, date)
		report.MixMelting = calc.MixMelting(pgdb, date)
		report.SiCC = calc.SiMnConsumption(pgdb, date)
		report.SiModel = calc.FeSiModelConsumption(pgdb, date)
		report.SiMnCC = calc.SiMnConsumption(pgdb, date)
		report.SiMnModel = calc.SiMnModelConsumption(pgdb, date)
		report.MnCC = calc.FeMnConsumption(pgdb, date)
		report.MnModel = calc.FeMnModelConsumption(pgdb, date)
		report.SlagTruncationRatio = calc.SlagTruncationRatio(pgdb, date)
		report.SlagSkimmingRatio = calc.SlagSkimmingRatio(pgdb, date)
		report.CCMeltingCycle = calc.CCMeltingCycleMinutes(pgdb, date)
		report.FePercentageInSlag = calc.FePercentageInSlag(pgdb, date)
		report.SlagSamplingPercentage = calc.SlagSamplingPercentage(pgdb, date)
		report.GoodCCOutput = calc.GoodCCOutput(pgdb, date)
		report.GoodCCMNLZOutput = calc.GoodCCMNLZOutput(pgdb, date)
		report.GoodIngotOutput = calc.GoodCCIngotOutput(pgdb, date)
		report.ProcessingTime = calc.ProcessingTime(pgdb, date)
		report.ArcTime = calc.ArcTime(pgdb, date)
		report.LimestoneConsumption = calc.LimestoneConsumption(pgdb, date)
		report.FluorsparConsumption = calc.FluorsparConsumption(pgdb, date)
		report.ArgonOxygenConsumption = calc.ArgonOxygenConsumption(pgdb, date)
		report.ElectricityConsumption = 0.0
		report.ElectrodeConsumption = 0.0
		report.InletTemperature = calc.InletTemperature(pgdb, date)
		report.InletOxidation = calc.InletOxidation(pgdb, date)

		database.InsertReport(msdb, *report)

		msdb.Close()
		pgdb.Close()
		logger.Info("Calculations is done!")
	}
}

func waitUntilMidnight() {
	currentTime := time.Now()
	targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 20, 0, 0, 0, currentTime.Location())
	timeToWait := targetTime.Sub(currentTime)
	logger.Info("The next calculation will be in", timeToWait)
	time.Sleep(timeToWait)
}
