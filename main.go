package main

import (
	"fmt"
	calc "main/calculations"
	c "main/configuration"
	"main/database"
	"main/logger"
	"main/models"
)

var report = new(models.Report)

func main() {
	logger.InitLogger()
	c.InitConfig("configuration/config.yaml")

	date := calc.GetDate()

	fmt.Println(date)

	msdev := database.ConnectMsDev()
	pgdev := database.ConnectPgDev()

	report.Date = "2023-07-14 01:00:00"
	report.CastIronMelting = calc.ConsumptionOfCastIronForMelting(pgdev, "2023-07-14 00:00:00")
	report.ScrapMelting = calc.ConsumptionOfScrapForMelting(pgdev, "2023-07-14 00:00:00")
	report.SiInCastIron = calc.GetSiInCastIron(pgdev, "2023-07-14 00:00:00")
	report.CastIronTemperature = calc.GetCastIronTemperature(pgdev, "2023-07-14 00:00:00")
	report.SContent = calc.GetSContent(pgdev, "2023-07-14 00:00:00")
	report.MNLZMelting = calc.MNLZMeltingAvgWeight(pgdev, "2023-07-14 00:00:00")
	report.IngotMelting = calc.IngotMeltingAvgWeight(pgdev, "2023-07-14 00:00:00")
	report.O2Content = calc.O2Content(pgdev, "2023-07-14 00:00:00")
	report.LimestoneFlow = calc.LimeFlow(pgdev, "2023-07-14 00:00:00")
	report.DolomiteFlow = calc.DolomiteFlow(pgdev, "2023-07-14 00:00:00")
	report.AluminumPreheating = calc.AluminumPreheating(pgdev, "2023-07-14 00:00:00")
	report.MixMelting = calc.MixMelting(pgdev, "2023-07-14 00:00:00")

	fmt.Println(report)
	database.InsertReport(msdev, *report)

	defer msdev.Close()
	defer pgdev.Close()
	logger.Info("Done!")
}
