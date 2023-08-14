package calculations

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/database"
	"main/models"

	"gonum.org/v1/gonum/floats"
)

// Потребление чугуна МНЛЗ
func castIronConsumptionMnlzSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ConsumptionMnlz, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Потребление чугуна слитки
func castIronConsumptionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ConsumptionIngot, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Количество плавок ОНРС
func numberMeltdownsOnrs(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.MeltdownsOnrs, date)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Количество плавок разливка
func numberMeltdownsCasting(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeat)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Потребление лома МНЛЗ
func scrapConsumptionMnlzSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ScrapConsumptionMnlz, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Потребление лома Слиток
func scrapConsumptionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ScrapConsumptionIngot, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Расход чугуна на плавку
func ConsumptionOfCastIronForMelting(db *sql.DB, date string) float64 {
	MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	BCConsumption := castIronConsumptionIngotSum(db, date)
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)

	ICCasting := SafeDivision((MNLZConsumption + BCConsumption), (SMeltCount + CPourCount))
	return ICCasting
}

// Расход лома на плавку
func ConsumptionOfScrapForMelting(db *sql.DB, date string) float64 {
	MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	BCScrapConsumption := scrapConsumptionIngotSum(db, date)
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)

	ScrapMelting := SafeDivision((MNLZScrapConsumption + BCScrapConsumption), (SMeltCount + CPourCount))
	return ScrapMelting
}

// Содержание Si в чугуне
func GetSiInCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SiIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Годный чугун, %
func GetGoodCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGoodCastIron, date, c.GlobalConfig.Measurings.Dtk)
	data := database.ExecuteQuery(db, q)
	countGoodCI := Len(data)

	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)

	res := SafeDivision(countGoodCI, (SMeltCount + CPourCount))
	return res
}

// Температура чугуна
func GetCastIronTemperature(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.TIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Содержание S, %
func GetSContent(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Производство МНЛЗ
func productionMNLZSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeatMnlz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Производство слитки
func productionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeat)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Средний вес плавки МНЛЗ, тонн
func MNLZMeltingAvgWeight(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	sMeltCount := numberMeltdownsOnrs(db, date)

	res := SafeDivision(prodMNLZ, sMeltCount)
	return res
}

// Средний вес плавки Слиток, тонн
func IngotMeltingAvgWeight(db *sql.DB, date string) float64 {
	prodIngot := productionIngotSum(db, date)
	cPourCount := numberMeltdownsCasting(db, date)

	res := SafeDivision(prodIngot, cPourCount)
	return res
}

// Содержание О2 на выпуске, ppm
func O2Content(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.OxiBof)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Известь вр.
func getLime(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WCaov)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Расход извести, кг/т
func LimeFlow(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	CCprod := prodMNLZ + prodIngot

	lime := getLime(db, date)
	res := SafeDivision(lime, (CCprod * 1000))
	return res
}

// Доломит
func getDolomite(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WDolo)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Расход доломита, кг/т
func DolomiteFlow(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	KCprod := prodMNLZ + prodIngot

	dolomite := getDolomite(db, date)
	res := SafeDivision(dolomite, (KCprod * 1000))
	return res
}

// Алюминий
func getAluminum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WAlrz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Алюминий на разогрев
func AluminumPreheating(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	CCprod := prodMNLZ + prodIngot

	alu := getAluminum(db, date)
	res := SafeDivision(alu, CCprod)
	return res
}

// Смесь
func getMixture(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WSmesrz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Смесь на плавку
func MixMelting(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	mix := getMixture(db, date)

	CCprod := prodMNLZ + prodIngot
	res := SafeDivision(mix, CCprod)
	return res
}

// Si КЦ всего
func getFeSiCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Si КЦ
func FeSiConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	fesisum := getFeSiCC(db, date)

	CCprod := prodMNLZ + prodIngot
	res := SafeDivision((fesisum * 1000), CCprod)
	return res
}

// FeSi65 расчет
func getSiModel(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Fesi65)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Si по модели
func FeSiModelConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	fesi65sum := getSiModel(db, date)

	res := SafeDivision(fesi65sum, prodMNLZ)
	return res
}

// SiMn КЦ всего
func getFeSiMnCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// SiMn КЦ
func SiMnConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	fesimnsum := getFeSiMnCC(db, date)

	CCprod := prodMNLZ + prodIngot
	res := SafeDivision((fesimnsum * 1000), CCprod)
	return res
}

// SiMn расчет
func getSiMnModel(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Simn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// SiMn по модели
func SiMnModelConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	simnmodelsum := getSiMnModel(db, date)

	res := SafeDivision(simnmodelsum, prodMNLZ)
	return res
}

// Mn КЦ всего
func getMnCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum1 := Sum(data)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn88)
	data = database.ExecuteQuery(db, q)
	sum2 := Sum(data)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn95)
	data = database.ExecuteQuery(db, q)
	sum3 := Sum(data)

	res := sum1 + sum2 + sum3
	return res
}

// Mn КЦ
func FeMnConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	mnCC := getMnCC(db, date)

	CCprod := prodMNLZ + prodIngot
	res := SafeDivision((mnCC * 1000), CCprod)
	return res
}

// Mn расчет
func getMnModel(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn78)
	data := database.ExecuteQuery(db, q)
	sum1 := Sum(data)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn88)
	data = database.ExecuteQuery(db, q)
	sum2 := Sum(data)

	return sum1 + sum2
}

// Mn по модели
func FeMnModelConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	mnmodel := getMnModel(db, date)

	res := SafeDivision(mnmodel, prodMNLZ)
	return res
}

// Признак отсечки Монокон
func getMonokonTruncationCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.FlHeatCutoffSlag)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Доля плавок с отсечкой шлака
func SlagTruncationRatio(db *sql.DB, date string) float64 {
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)
	MonokonCount := getMonokonTruncationCount(db, date)

	res := SafeDivision(MonokonCount, (SMeltCount + CPourCount))
	return res
}

// Признак
func getSkimmingFlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.FlSlop)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Доля плавок со скачиванием шлака
func SlagSkimmingRatio(db *sql.DB, date string) float64 {
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)
	SlagSkimmingCount := getSkimmingFlagCount(db, date)

	res := SafeDivision(SlagSkimmingCount, (SMeltCount + CPourCount))
	return res
}

// dtn
func getDtn(db *sql.DB, date string) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Dtn)
	data := database.ExecuteQuery(db, q)
	return data
}

// dtk
func getDtk(db *sql.DB, date string) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Dtk)
	data := database.ExecuteQuery(db, q)
	return data
}

// Цикл плавки КЦ, мин
func CCMeltingCycleMinutes(db *sql.DB, date string) float64 {
	dtn := getDtn(db, date)
	dtk := getDtk(db, date)

	res := AvgDiffDate(dtn, dtk)
	return res
}

// % Fe в шлаке
func FePercentageInSlag(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.ShlMgo)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

func getFePercentageInSlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.ShlMgo)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// % отбора проб шлака
func SlagSamplingPercentage(db *sql.DB, date string) float64 {
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)
	SlagCount := getFePercentageInSlagCount(db, date)

	res := SafeDivision(SlagCount, (SMeltCount+CPourCount)) * 100
	return res
}

// w_FeMn78,  w_heat_mnlz > 0
func getWFeMn78MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeMn88, w_heat_mnlz > 0
func getWFeMn88MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn88)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeMn95, w_heat_mnlz > 0
func getWFeMn95MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn95)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeSi, w_heat_mnlz > 0
func getWFeSiMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_SiMn, w_heat_mnlz > 0
func getWSiMnMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr, w_heat_mnlz > 0
func getWFeCrMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr025, w_heat_mnlz > 0
func getWFeCr025MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr025)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr800, w_heat_mnlz > 0
func getWFeCr800MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr800)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// всего феросплавов для МНЛЗ
func ferroalloysOnMNLZ(db *sql.DB, date string) float64 {
	femn78 := getWFeMn78MNLZ(db, date)
	femn88 := getWFeMn88MNLZ(db, date)
	femn95 := getWFeMn95MNLZ(db, date)
	fesi := getWFeSiMNLZ(db, date)
	simn := getWSiMnMNLZ(db, date)
	fecr := getWFeCrMNLZ(db, date)
	fecr25 := getWFeCr025MNLZ(db, date)
	fecr800 := getWFeCr800MNLZ(db, date)

	res := femn78 + femn88 + femn95 + fesi + simn + fecr + fecr25 + fecr800
	return res
}

// w_FeMn78, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn78Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeMn88, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn88Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn88)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeMn95, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn95Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn95)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeSi, w_iron > 0 && w_heat_mnlz > 0
func getWFeSiIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_SiMn, w_iron > 0 && w_heat_mnlz > 0
func getWSiMnIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr, w_iron > 0 && w_heat_mnlz > 0
func getWFeCrIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr025, w_iron > 0 && w_heat_mnlz > 0
func getWFeCr025Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr025)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// w_FeCr800, w_iron > 0 && w_heat_mnlz > 0
func getWFeCr800Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr800)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// всего феросплавов для Слитков
func ferroalloysOnIngot(db *sql.DB, date string) float64 {
	femn78 := getWFeMn78Ingot(db, date)
	femn88 := getWFeMn88Ingot(db, date)
	femn95 := getWFeMn95Ingot(db, date)
	fesi := getWFeSiIngot(db, date)
	simn := getWSiMnIngot(db, date)
	fecr := getWFeCrIngot(db, date)
	fecr25 := getWFeCr025Ingot(db, date)
	fecr800 := getWFeCr800Ingot(db, date)

	res := femn78 + femn88 + femn95 + fesi + simn + fecr + fecr25 + fecr800
	return res
}

// Выход годного КЦ общий
func GoodCCOutput(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	prodCC := prodMNLZ + prodIngot

	MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	BCConsumption := castIronConsumptionIngotSum(db, date)
	CCIronConsumption := MNLZConsumption + BCConsumption

	ferroalloysMNLZ := ferroalloysOnMNLZ(db, date)
	ferroalloysIngot := ferroalloysOnIngot(db, date)
	ferroalloys := ferroalloysIngot + ferroalloysMNLZ

	MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	BCScrapConsumption := scrapConsumptionIngotSum(db, date)

	res := SafeDivision(prodCC, (CCIronConsumption + MNLZScrapConsumption + BCScrapConsumption + ferroalloys/1000))
	return res
}

// Выход годного КЦ МНЛЗ
func GoodCCMNLZOutput(db *sql.DB, date string) float64 {
	MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	ferroalloysMNLZ := ferroalloysOnMNLZ(db, date)

	res := SafeDivision(MNLZConsumption, (MNLZScrapConsumption + ferroalloysMNLZ/1000))
	return res
}

// Выход годного КЦ Слиток
func GoodCCIngotOutput(db *sql.DB, date string) float64 {
	prodIngot := productionIngotSum(db, date)
	BCConsumption := castIronConsumptionIngotSum(db, date)
	ferroalloysIngot := ferroalloysOnIngot(db, date)
	BCScrapConsumption := scrapConsumptionIngotSum(db, date)

	res := SafeDivision(prodIngot, (BCConsumption + BCScrapConsumption + ferroalloysIngot))
	return res
}

// Получить УПК за номером
func getUPKByNumber(db *sql.DB, date string, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetUpk, date, c.GlobalConfig.Measurings.TmTreatment, n)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Время обработки
func ProcessingTime(db *sql.DB, date string) float64 {
	res := 0.0
	for n := 1; n <= 3; n++ {
		pTime := getUPKByNumber(db, date, n)
		res = res + pTime
	}
	res = res / 3
	return res
}

// Получить нагрев за номерром УПК
func getHeatingByUPKNumber(db *sql.DB, date string, id int, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetUpk, date, id, n)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Время дуги
func ArcTime(db *sql.DB, date string) float64 {
	res := []float64{0.0, 0.0, 0.0}
	tmp := []int{c.GlobalConfig.Measurings.M1e,
		c.GlobalConfig.Measurings.M2e,
		c.GlobalConfig.Measurings.M3e,
		c.GlobalConfig.Measurings.M4e,
		c.GlobalConfig.Measurings.M5e,
		c.GlobalConfig.Measurings.M6e,
		c.GlobalConfig.Measurings.M7e,
		c.GlobalConfig.Measurings.M8e}
	for n := 1; n <= 3; n++ {
		for _, id := range tmp {
			h := getHeatingByUPKNumber(db, date, id, n)
			res[n-1] = res[n-1] + h
		}
		res[n-1] = res[n-1] * 24
	}

	calc := floats.Sum(res)
	calc = calc / 3
	return calc
}

// Известь
func getCaO(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCaoLf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Известь высоко-кальц
func getCaOHighCalcium(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCao90Lf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Потребление Извести, кг/т
func LimestoneConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	CaO := getCaO(db, date)
	CaOhigh := getCaOHighCalcium(db, date)

	res := SafeDivision((CaO + CaOhigh), prodMNLZ)
	return res
}

// Шпат
func getCaF2(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCaf2Lf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Потребление шпата
func FluorsparConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	caF2 := getCaF2(db, date)

	res := SafeDivision(caF2, prodMNLZ)
	return res
}

// АРК
func getAPK(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WApkLf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Потребление АРК
func ArgonOxygenConsumption(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	apk := getAPK(db, date)

	res := SafeDivision(apk, prodMNLZ)
	return res
}
