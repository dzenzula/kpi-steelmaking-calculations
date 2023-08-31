package main

import (
	"fmt"
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
		//waitUntilMidnight()

		date := calc.GetDate()
		msdev := database.ConnectMsDev()
		pgdev := database.ConnectPgDev()

		report.Date = date
		report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdev, date)
		report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdev, date)
		report.SiInCastIron = calc.GetSiInCastIron(pgdev, date)
		report.CastIronTemperature = calc.GetCastIronTemperature(pgdev, date)
		report.GoodCastIron = calc.GetGoodCastIron(pgdev, date)
		report.SContent = calc.GetSContent(pgdev, date)
		report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdev, date)
		report.IngotMelting = calc.IngotMeltingAvgWeight(pgdev, date)
		report.O2Content = calc.O2Content(pgdev, date)
		report.LimestoneFlow = calc.LimeFlow(pgdev, date)
		report.DolomiteFlow = calc.DolomiteFlow(pgdev, date)
		report.AluminumPreheating = calc.AluminumPreheating(pgdev, date)
		report.MixMelting = calc.MixMelting(pgdev, date)
		report.SiCC = calc.SiMnConsumption(pgdev, date)
		report.SiModel = calc.FeSiModelConsumption(pgdev, date)
		report.SiMnCC = calc.SiMnConsumption(pgdev, date)
		report.SiMnModel = calc.SiMnModelConsumption(pgdev, date)
		report.MnCC = calc.FeMnConsumption(pgdev, date)
		report.MnModel = calc.FeMnModelConsumption(pgdev, date)
		report.SlagTruncationRatio = calc.SlagTruncationRatio(pgdev, date)
		report.SlagSkimmingRatio = calc.SlagSkimmingRatio(pgdev, date)
		report.CCMeltingCycle = calc.CCMeltingCycleMinutes(pgdev, date)
		report.FePercentageInSlag = calc.FePercentageInSlag(pgdev, date)
		report.SlagSamplingPercentage = calc.SlagSamplingPercentage(pgdev, date)
		report.GoodCCOutput = calc.GoodCCOutput(pgdev, date)
		report.GoodCCMNLZOutput = calc.GoodCCMNLZOutput(pgdev, date)
		report.GoodIngotOutput = calc.GoodCCIngotOutput(pgdev, date)
		report.ProcessingTime = calc.ProcessingTime(pgdev, date)
		report.ArcTime = calc.ArcTime(pgdev, date)
		report.LimestoneConsumption = calc.LimestoneConsumption(pgdev, date)
		report.FluorsparConsumption = calc.FluorsparConsumption(pgdev, date)
		report.ArgonOxygenConsumption = calc.ArgonOxygenConsumption(pgdev, date)
		report.ElectricityConsumption = 0.0
		report.ElectrodeConsumption = 0.0
		report.InletTemperature = calc.InletTemperature(pgdev, date)
		report.InletOxidation = calc.InletOxidation(pgdev, date)

		fmt.Println(date)
		database.InsertReport(msdev, *report)

		msdev.Close()
		pgdev.Close()
		logger.Info("Done!")
	}
}

func waitUntilMidnight() {
	currentTime := time.Now()
	targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 20, 0, 0, 0, currentTime.Location()).Add(24 * time.Hour)
	timeToWait := targetTime.Sub(currentTime)
	fmt.Println(timeToWait)
	time.Sleep(timeToWait)
}
