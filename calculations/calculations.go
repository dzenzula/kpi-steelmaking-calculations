package calculations

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/database"
	"main/models"
	"strconv"

	"gonum.org/v1/gonum/floats"
)

var (
	mnlzConsumption      float64
	bcConsumption        float64
	sMeltCount           float64
	cPourCount           float64
	mnlzScrapConsumption float64
	bcScrapConsumption   float64
	prodMNLZ             float64
	prodIngot            float64
	ferroalloysMNLZ      float64
	ferroalloysIngot     float64
)

func CacheInit(db *sql.DB, date string) {
	mnlzConsumption = castIronConsumptionMnlzSum(db, date)
	bcConsumption = castIronConsumptionIngotSum(db, date)
	sMeltCount = numberMeltdownsOnrs(db, date)
	cPourCount = numberMeltdownsCasting(db, date)
	mnlzScrapConsumption = scrapConsumptionMnlzSum(db, date)
	bcScrapConsumption = scrapConsumptionIngotSum(db, date)
	prodMNLZ = productionMNLZSum(db, date)
	prodIngot = productionIngotSum(db, date)
	ferroalloysMNLZ = ferroalloysOnMNLZ(db, date)
	ferroalloysIngot = ferroalloysOnIngot(db, date)
}

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
	//MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	//BCConsumption := castIronConsumptionIngotSum(db, date)
	//SMeltCount := numberMeltdownsOnrs(db, date)
	//CPourCount := numberMeltdownsCasting(db, date)

	ICCasting := SafeDivision((mnlzConsumption + bcConsumption), (sMeltCount + cPourCount))

	fmt.Println("Расход чугуна = Потребление чугуна КЦ / Плавок по КЦ")
	fmt.Printf("%f = %f / %f \n", ICCasting, (mnlzConsumption + bcConsumption), (sMeltCount + cPourCount))
	return ICCasting
}

// Расход лома на плавку
func ConsumptionOfScrapForMelting(db *sql.DB, date string) float64 {
	//MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	//BCScrapConsumption := scrapConsumptionIngotSum(db, date)
	//SMeltCount := numberMeltdownsOnrs(db, date)
	//CPourCount := numberMeltdownsCasting(db, date)

	ScrapMelting := SafeDivision((mnlzScrapConsumption + bcScrapConsumption), (sMeltCount + cPourCount))
	fmt.Println("Расход лома = Потребление чугуна КЦ / Плавок по КЦ")
	fmt.Printf("%f = %f / %f \n", ScrapMelting, (mnlzScrapConsumption + bcScrapConsumption), (sMeltCount + cPourCount))
	return ScrapMelting
}

// Содержание Si в чугуне
func GetSiInCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SiIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	fmt.Println("Si в чугуне =", avg)
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
	fmt.Println("Годного чугуна = Годного чугуна - кол-во записей / (Количество плавок OHPC + Количество плавок разливка )")
	fmt.Printf("%f = %f / %f \n", res, countGoodCI, (SMeltCount + CPourCount))
	return res
}

// Температура чугуна
func GetCastIronTemperature(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.TIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	fmt.Println("Температура чугуна =", avg)
	return avg
}

// Содержание S, %
func GetSContent(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	fmt.Println("Содержание S =", avg)
	return avg
}

// Производство МНЛЗ,
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
	fmt.Println("Средний вес плавки MNLZ = Производство МНЛЗ / Количество плавок OHPC")
	fmt.Printf("%f = %f / %f \n", res, prodMNLZ, sMeltCount)
	return res
}

// Средний вес плавки Слиток, тонн
func IngotMeltingAvgWeight(db *sql.DB, date string) float64 {
	//prodIngot := productionIngotSum(db, date)
	//cPourCount := numberMeltdownsCasting(db, date)

	res := SafeDivision(prodIngot, cPourCount)
	fmt.Println("Средний вес плавки Слиток = Производство слитки / Количество плавок разливка")
	fmt.Printf("%f = %f / %f \n", res, prodIngot, cPourCount)
	return res
}

// Содержание О2 на выпуске, ppm
func O2Content(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.OxiBof)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	fmt.Println("Содержание О2 на выпуске =", avg)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	ccProd := prodMNLZ + prodIngot

	lime := getLime(db, date)
	res := SafeDivision(lime, (ccProd * 1000))
	fmt.Println("Известь на плавку = Известь вр. / (Производство КЦ * 1000)")
	fmt.Printf("%f = %f / (%f * 1000) \n", res, lime, ccProd)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	kcProd := prodMNLZ + prodIngot

	dolomite := getDolomite(db, date)
	res := SafeDivision(dolomite, (kcProd * 1000))
	fmt.Println("Доломит на плавку = Доло-мит / (Производство КЦ * 1000)")
	fmt.Printf("%f = %f / (%f * 1000) \n", res, dolomite, kcProd)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	ccProd := prodMNLZ + prodIngot

	alu := getAluminum(db, date)
	res := SafeDivision(alu, ccProd)
	fmt.Println("Алюминий на разогрев = Алюминий / Производство КЦ")
	fmt.Printf("%f = %f / %f\n", res, alu, ccProd)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	mix := getMixture(db, date)

	cCprod := prodMNLZ + prodIngot
	res := SafeDivision(mix, cCprod)
	fmt.Println("Смесь на плавку = Смесь / Производство КЦ")
	fmt.Printf("%f = %f / %f\n", res, mix, cCprod)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	fesisum := getFeSiCC(db, date)

	cCprod := prodMNLZ + prodIngot
	res := SafeDivision((fesisum * 1000), cCprod)
	fmt.Println("Si КЦ = Si КЦ всего * 1000 / Производство КЦ")
	fmt.Printf("%f = %f * 1000 / %f\n", res, fesisum, cCprod)
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
	//prodMNLZ := productionMNLZSum(db, date)
	fesi65sum := getSiModel(db, date)

	res := SafeDivision(fesi65sum, prodMNLZ)
	fmt.Println("Si по модели = FeSi65 расчет / Производство МНЛЗ")
	fmt.Printf("%f = %f / %f\n", res, fesi65sum, prodMNLZ)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	fesimnsum := getFeSiMnCC(db, date)

	cCprod := prodMNLZ + prodIngot
	res := SafeDivision((fesimnsum * 1000), cCprod)
	fmt.Println("SiMn КЦ = SiMn КЦ всего * 1000 / Производство КЦ")
	fmt.Printf("%f = %f * 1000 / %f\n", res, fesimnsum, cCprod)
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
	//prodMNLZ := productionMNLZSum(db, date)
	simnmodelsum := getSiMnModel(db, date)

	res := SafeDivision(simnmodelsum, prodMNLZ)
	fmt.Println("SiMn по модели = SiMn расчет / Производство МНЛЗ")
	fmt.Printf("%f = %f / %f\n", res, simnmodelsum, prodMNLZ)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	mnCC := getMnCC(db, date)

	cCprod := prodMNLZ + prodIngot
	res := SafeDivision((mnCC * 1000), cCprod)
	fmt.Println("Mn КЦ = Mn КЦ всего * 1000 / Производство КЦ")
	fmt.Printf("%f = %f * 1000 / %f\n", res, mnCC, cCprod)
	return res
}

// Mn расчет
func getMnModel(db *sql.DB, date string) (float64, float64) {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn78)
	data := database.ExecuteQuery(db, q)
	sum1 := Sum(data)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn88)
	data = database.ExecuteQuery(db, q)
	sum2 := Sum(data)

	return sum1, sum2
}

// Mn по модели
func FeMnModelConsumption(db *sql.DB, date string) float64 {
	//prodMNLZ := productionMNLZSum(db, date)
	feMn78, femn88 := getMnModel(db, date)
	mnmodel := feMn78 + femn88
	res := SafeDivision(mnmodel, prodMNLZ)
	fmt.Println("Mn по модели = (FeMn78 Расчет + FeMn88 Расчет) / Производство МНЛЗ")
	fmt.Printf("%f = (%f + %f) / %f\n", res, feMn78, femn88, prodMNLZ)
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
	//SMeltCount := numberMeltdownsOnrs(db, date)
	//CPourCount := numberMeltdownsCasting(db, date)
	MonokonCount := getMonokonTruncationCount(db, date)

	res := SafeDivision(MonokonCount, (sMeltCount + cPourCount))
	fmt.Println("Отсечка шлака = Количество признаков отсечки Монокон / (Количество плавок OHPC + Количество плавок разливка)")
	fmt.Printf("%f = %f / (%f + %f)\n", res, MonokonCount, sMeltCount, cPourCount)
	return res
}

// Признак
func getSkimmingFlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetSlagTruncation, date)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Доля плавок со скачиванием шлака
func SlagSkimmingRatio(db *sql.DB, date string) float64 {
	//SMeltCount := numberMeltdownsOnrs(db, date)
	//CPourCount := numberMeltdownsCasting(db, date)
	SlagSkimmingCount := getSkimmingFlagCount(db, date)

	res := SafeDivision(SlagSkimmingCount, (sMeltCount + cPourCount))
	fmt.Println("Скачивание шлака = Количество признаков скачки / (Количество плавок OHPC + Количество плавок разливка)")
	fmt.Printf("%f = %f / (%f + %f)\n", res, SlagSkimmingCount, sMeltCount, cPourCount)
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
	fmt.Println("Длительность плавки =  конец плавки - начало плавки")
	fmt.Printf("Средняя длительность плавки = %f\n", res)
	return res
}

// % Fe в шлаке
func FePercentageInSlag(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.ShlMgo)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	fmt.Printf("Fe в шлаке = %f\n", avg)
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
	//SMeltCount := numberMeltdownsOnrs(db, date)
	//CPourCount := numberMeltdownsCasting(db, date)
	slagCount := getFePercentageInSlagCount(db, date)

	res := SafeDivision(slagCount, (sMeltCount+cPourCount)) * 100
	fmt.Println("Кол-во шлаков = MgO / (Количество плавок OHPC + Количество плавок разливка) * 100")
	fmt.Printf("%f = %f / (%f + %f) * 100\n", res, slagCount, sMeltCount, cPourCount)
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
	//prodMNLZ := productionMNLZSum(db, date)
	//prodIngot := productionIngotSum(db, date)
	prodCC := prodMNLZ + prodIngot

	//MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	//BCConsumption := castIronConsumptionIngotSum(db, date)
	CCIronConsumption := mnlzConsumption + bcConsumption

	//ferroalloysMNLZ := ferroalloysOnMNLZ(db, date)
	//ferroalloysIngot := ferroalloysOnIngot(db, date)
	ferroalloys := ferroalloysIngot + ferroalloysMNLZ

	//MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	//BCScrapConsumption := scrapConsumptionIngotSum(db, date)

	res := SafeDivision(prodCC, (CCIronConsumption + mnlzScrapConsumption + bcScrapConsumption + ferroalloys/1000))
	fmt.Println("Выход годного КЦ = Производство КЦ / (Потребление чугуна КЦ + Потребление лома МНЛЗ + Потребление лома Слиток + всего феросплавов КЦ /1000)")
	fmt.Printf("%f = %f / (%f + %f + %f + %f / 1000)\n", res, prodCC, CCIronConsumption, mnlzScrapConsumption, bcScrapConsumption, ferroalloys)
	return res
}

// Выход годного КЦ МНЛЗ
func GoodCCMNLZOutput(db *sql.DB, date string) float64 {
	//MNLZConsumption := castIronConsumptionMnlzSum(db, date)
	//MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	//ferroalloysMNLZ := ferroalloysOnMNLZ(db, date)

	res := SafeDivision(mnlzConsumption, (mnlzScrapConsumption + ferroalloysMNLZ/1000))
	fmt.Println("Выход годного КЦ МНЛЗ = Потребление чугуна МНЛЗ / (Потребление лома МНЛЗ + всего феросплавов для МНЛЗ / 1000)")
	fmt.Printf("%f = %f / (%f + %f / 1000)\n", res, mnlzConsumption, mnlzScrapConsumption, ferroalloysMNLZ)
	return res
}

// Выход годного КЦ Слиток
func GoodCCIngotOutput(db *sql.DB, date string) float64 {
	//prodIngot := productionIngotSum(db, date)
	//BCConsumption := castIronConsumptionIngotSum(db, date)
	//ferroalloysIngot := ferroalloysOnIngot(db, date)
	//BCScrapConsumption := scrapConsumptionIngotSum(db, date)

	res := SafeDivision(prodIngot, (bcConsumption + bcScrapConsumption + ferroalloysIngot))
	fmt.Println("Выход годного Слиток = Производство слитки / (Потребление чугуна слитки + Потребление лома Слиток + всего феросплавов для слитка))")
	fmt.Printf("%f = %f / (%f + %f + %f)\n", res, prodIngot, bcConsumption, bcScrapConsumption, ferroalloysIngot)
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
	fmt.Println("Время обработки = (УПК1 + УПК2 + УПК3) / 3")
	fmt.Printf("Время обработки = %f\n", res)
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
	fmt.Printf("Время дуги = %f\n", calc)
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
	//prodMNLZ := productionMNLZSum(db, date)
	CaO := getCaO(db, date)
	CaOhigh := getCaOHighCalcium(db, date)

	res := SafeDivision((CaO + CaOhigh), prodMNLZ)
	fmt.Println("Известь = Известь / Производство МНЛЗ")
	fmt.Printf("%f = %f / %f\n", res, (CaO + CaOhigh), prodMNLZ)
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
	//prodMNLZ := productionMNLZSum(db, date)
	caF2 := getCaF2(db, date)

	res := SafeDivision(caF2, prodMNLZ)
	fmt.Println("Шпат = Шпат / Производство МНЛЗ")
	fmt.Printf("%f = %f / %f\n", caF2, caF2, prodMNLZ)
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
	//prodMNLZ := productionMNLZSum(db, date)
	apk := getAPK(db, date)

	res := SafeDivision(apk, prodMNLZ)
	fmt.Println("АРК = АРК / Производство МНЛЗ")
	fmt.Printf("%f = %f / %f\n", apk, apk, prodMNLZ)
	return res
}

// Расходу Электричества
func ElectricityConsumption(db *sql.DB, date string) float64 {
	return 0.0
}

// Расход Электродов
func ElectrodeConsumption(db *sql.DB, date string) float64 {
	return 0.0
}

// Температура по приходу
func InletTemperature(db *sql.DB, date string) float64 {
	res := 0.0
	for i := 1; i < 4; i++ {
		q := fmt.Sprintf(c.GlobalConfig.Querries.GetInletTemperatureOxidation, date, c.GlobalConfig.Measurings.T1, i)
		data := database.ExecuteQuery(db, q)
		avg := Avg(data)
		res = res + avg
	}

	return res / 3
}

// Окисленность по приходу
func InletOxidation(db *sql.DB, date string) float64 {
	res := 0.0
	for i := 1; i < 4; i++ {
		q := fmt.Sprintf(c.GlobalConfig.Querries.GetInletTemperatureOxidation, date, c.GlobalConfig.Measurings.O21, i)
		data := database.ExecuteQuery(db, q)
		avg := Avg(data)
		res = res + avg
	}
	return res / 3
}

// Количество шлаков
func getSlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SampleTime)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Количество плавок
func getMeltingCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.HeatStart)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	return len
}

// Анализ шлаков УПК
func UPKSlagAnalysis(db *sql.DB, date string) float64 {
	slagCount := getSlagCount(db, date)
	meltCount := getMeltingCount(db, date)

	res := SafeDivision(slagCount, meltCount)
	fmt.Println("Анализ шлаков УПК = количество шлаков / количество плавок")
	fmt.Printf("%f = %f / %f\n", res, slagCount, meltCount)
	return res
}

// МНЛЗ Открытие
func getOpening(db *sql.DB, date string, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.DtBegin, i)
	data := database.ExecuteQuery(db, q)
	return data
}

// МНЛЗ Закрытие
func getClosing(db *sql.DB, date string, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.DtEnd, i)
	data := database.ExecuteQuery(db, q)
	return data
}

// Цикл разливки
func CastingCycle(db *sql.DB, date string) float64 {
	var dtn, dtk []models.Query
	for i := 1; i <= 3; i++ {
		dtn = append(dtn, getOpening(db, date, i)...)
		dtk = append(dtk, getClosing(db, date, i)...)
	}
	res := AvgDiffDate(dtn, dtk)
	fmt.Println("Цикл разливки =  Закрытие - Открытие")
	fmt.Printf("Цикл разливки = %f\n", res)
	return res
}

// Ср. скорость ручья
func getAvgSpeed(db *sql.DB, date string, i int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.AvgSpeed, i)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Скорость разливки
func CastingSpeed(db *sql.DB, date string) float64 {
	res := 0.0
	for i := 1; i <= 3; i++ {
		speed := getAvgSpeed(db, date, i)
		res = res + speed
	}
	res = res / 3
	return res
}

// Расход смазки
func getGreaseConsumption(db *sql.DB, date string, id int, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, id, i)
	data := database.ExecuteQuery(db, q)
	return data
}

// Конец серии
func getEndSeries(db *sql.DB, date string, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.EndSeries, i)
	data := database.ExecuteQuery(db, q)
	return data
}

// Серийность стопорной разливки
func CastingStopperSerial(db *sql.DB, date string) float64 {
	tmp := []int{
		c.GlobalConfig.Measurings.S1Oil,
		c.GlobalConfig.Measurings.S2Oil,
		c.GlobalConfig.Measurings.S3Oil,
		c.GlobalConfig.Measurings.S4Oil,
		c.GlobalConfig.Measurings.S5Oil,
		c.GlobalConfig.Measurings.S6Oil,
	}

	for n := 1; n <= 3; n++ {
		averages := make([]float64, len(tmp))
		valueArrays := make([][]float64, len(tmp))

		for i, id := range tmp {
			data := getGreaseConsumption(db, date, id, n)

			values := make([]float64, len(data))

			for j, query := range data {
				if query.Value != nil {
					value, err := strconv.ParseFloat(*query.Value, 64)
					if err == nil {
						values[j] = value
					}
				}
			}

			valueArrays[i] = values
		}

		for i, values := range valueArrays {
			if len(values) > 0 {
				sum := 0.0
				for _, value := range values {
					sum += value
				}
				average := sum / float64(len(values))
				averages[i] = average
			}
		}
	}

	return 0.0
}

// Серийность открытой разливки МНЛЗ1
func MNLZ1OpenSerial(db *sql.DB, date string) float64 {
	return 0.0
}

// Серийность открытой разливки МНЛЗ2
func MNLZ2OpenSerial(db *sql.DB, date string) float64 {
	return 0.0
}

// Серийность открытой разливки МНЛЗ3
func MNLZ3OpenSerial(db *sql.DB, date string) float64 {
	return 0.0
}

// Кол-во ручьев
func getStreamsCount(db *sql.DB, date string, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.PCount, n)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Количество ручьев в работе МНЛЗ1
func MNLZ1Streams(db *sql.DB, date string) float64 {
	res := getStreamsCount(db, date, 1)
	return res
}

// Количество ручьев в работе МНЛЗ2
func MNLZ2Streams(db *sql.DB, date string) float64 {
	res := getStreamsCount(db, date, 2)
	return res
}

// Количество ручьев в работе МНЛЗ3
func MNLZ3Streams(db *sql.DB, date string) float64 {
	res := getStreamsCount(db, date, 3)
	return res
}

//

// Длительность перепаковки МНЛЗ1, мин
func MNLZ1RepackingDuration(db *sql.DB, date string) float64 {
	return 0.0
}

// Длительность перепаковки МНЛЗ1, мин
func MNLZ2RepackingDuration(db *sql.DB, date string) float64 {
	return 0.0
}

// Длительность перепаковки МНЛЗ1, мин
func MNLZ3RepackingDuration(db *sql.DB, date string) float64 {
	return 0.0
}

// Плавки с отклонением по температуре МНЛЗ1, %
func MNLZ1MeltTempDeviation(db *sql.DB, date string) float64 {
	return 0.0
}

// Плавки с отклонением по температуре МНЛЗ2, %
func MNLZ2MeltTempDeviation(db *sql.DB, date string) float64 {
	return 0.0
}

// Плавки с отклонением по температуре МНЛЗ3, %
func MNLZ3MeltTempDeviation(db *sql.DB, date string) float64 {
	return 0.0
}

// Выход годного МНЛЗ
func GoodMNLZOutput(db *sql.DB, date string) float64 {
	return 0.0
}

// Время нахождения меиалла в ковше (до разливки), мин
func MetalRetentionTime(db *sql.DB, date string) float64 {
	return 0.0
}
