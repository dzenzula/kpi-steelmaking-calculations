package calculations

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/database"
	"main/logger"
	"main/models"
	"math"
	"strconv"
	"strings"

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
	idOilList            []int
)

func CacheInit(db *sql.DB, date string) {
	logger.Info("Calculation cache initialization started")
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
	idOilList = []int{
		c.GlobalConfig.Measurings.S1Oil,
		c.GlobalConfig.Measurings.S2Oil,
		c.GlobalConfig.Measurings.S3Oil,
		c.GlobalConfig.Measurings.S4Oil,
		c.GlobalConfig.Measurings.S5Oil,
		c.GlobalConfig.Measurings.S6Oil,
	}
	logger.Info("Calculation cache initialization ended")
}

// Потребление чугуна МНЛЗ
func castIronConsumptionMnlzSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ConsumptionMnlz, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	logger.Debug("Потребление чугуна МНЛЗ =", sum)
	return sum
}

// Потребление чугуна слитки
func castIronConsumptionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ConsumptionIngot, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	logger.Debug("Потребление чугуна слитки = ", sum)
	return sum
}

// Количество плавок ОНРС
func numberMeltdownsOnrs(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.MeltdownsOnrs, date)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	logger.Debug("Количество плавок ОНРС = ", len)
	return len
}

// Количество плавок разливка
func numberMeltdownsCasting(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeat)
	data := database.ExecuteQuery(db, q)
	len := Len(data)
	logger.Debug("Количество плавок разливка = ", len)
	return len
}

// Потребление лома МНЛЗ
func scrapConsumptionMnlzSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ScrapConsumptionMnlz, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	logger.Debug("Потребление лома МНЛЗ = ", sum)
	return sum
}

// Потребление лома Слиток
func scrapConsumptionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.ScrapConsumptionIngot, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	logger.Debug("Потребление лома Слиток = ", sum)
	return sum
}

// Расход чугуна на плавку
func ConsumptionOfCastIronForMelting(db *sql.DB, date string) float64 {
	ICCasting := SafeDivision((mnlzConsumption + bcConsumption), (sMeltCount + cPourCount))

	logger.Debug("Потребление чугуна КЦ = Потребление чугуна МНЛЗ + Потребление чугуна слитки")
	logger.Debug("Потребление чугуна КЦ = ", (mnlzConsumption + bcConsumption))
	logger.Debug("Плавок по КЦ = Количество плавок ОНРС + Количество плавок разливка")
	logger.Debug("Плавок по КЦ = ", (sMeltCount + cPourCount))
	logger.Debug("Расход чугуна = Потребление чугуна КЦ / Плавок по КЦ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", ICCasting, (mnlzConsumption + bcConsumption), (sMeltCount + cPourCount)))
	return ICCasting
}

// Расход лома на плавку
func ConsumptionOfScrapForMelting(db *sql.DB, date string) float64 {
	ScrapMelting := SafeDivision((mnlzScrapConsumption + bcScrapConsumption), (sMeltCount + cPourCount))

	logger.Debug("Потребление чугуна КЦ = Потребление лома МНЛЗ + Потребление лома слитки")
	logger.Debug("Потребление чугуна КЦ =", (mnlzScrapConsumption + bcScrapConsumption))
	logger.Debug("Расход лома = Потребление лома КЦ / Плавок по КЦ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", ScrapMelting, (mnlzScrapConsumption + bcScrapConsumption), (sMeltCount + cPourCount)))
	return ScrapMelting
}

// Содержание Si в чугуне
func GetSiInCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SiIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug("Si в чугуне = ", avg)
	return avg
}

// Годный чугун, %
func GetGoodCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGoodCastIron, date, c.GlobalConfig.Measurings.Dtk)
	data := database.ExecuteQuery(db, q)
	countGoodCI := Len(data)
	res := SafeDivision(countGoodCI, (sMeltCount + cPourCount))

	logger.Debug("Годного чугуна - кол-во записей = ", countGoodCI)
	logger.Debug("Плавок по КЦ = Количество плавок ОНРС + Количество плавок разливка")
	logger.Debug("Плавок по КЦ = ", (sMeltCount + cPourCount))
	logger.Debug("Годного чугуна = Годного чугуна - кол-во записей / Плавок по КЦ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, countGoodCI, (sMeltCount + cPourCount)))
	return res
}

// Температура чугуна
func GetCastIronTemperature(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.TIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug("Температура чугуна = ", avg)
	return avg
}

// Содержание S, %
func GetSContent(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SIron)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	logger.Debug("Содержание S = ", avg)
	return avg
}

// Производство МНЛЗ
func productionMNLZSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeatMnlz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Производство МНЛЗ = ", sum)
	return sum
}

// Производство слитки
func productionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WHeat)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Производство слитки = ", sum)
	return sum
}

// Средний вес плавки МНЛЗ, тонн
func MNLZMeltingAvgWeight(db *sql.DB, date string) float64 {
	res := SafeDivision(prodMNLZ, sMeltCount)
	logger.Debug("Средний вес плавки MNLZ = Производство МНЛЗ / Количество плавок OHPC")
	logger.Debug(fmt.Sprintf("%f = %f / %f ", res, prodMNLZ, sMeltCount))
	return res
}

// Средний вес плавки Слиток, тонн
func IngotMeltingAvgWeight(db *sql.DB, date string) float64 {
	res := SafeDivision(prodIngot, cPourCount)
	logger.Debug("Средний вес плавки Слиток = Производство слитки / Количество плавок разливка")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, prodIngot, cPourCount))
	return res
}

// Содержание О2 на выпуске, ppm
func O2Content(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.OxiBof)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug("Содержание О2 на выпуске = ", avg)
	return avg
}

// Известь вр.
func getLime(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WCaov)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Известь = ", sum)
	return sum
}

// Расход извести, кг/т
func LimeFlow(db *sql.DB, date string) float64 {
	ccProd := prodMNLZ + prodIngot
	lime := getLime(db, date)
	res := SafeDivision(lime, (ccProd * 1000))

	logger.Debug("Производство КЦ = Производство МНЛЗ + Производство слитки")
	logger.Debug("Производство КЦ = ", ccProd)
	logger.Debug("Известь на плавку = Известь / (Производство КЦ * 1000)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f * 1000)", res, lime, ccProd))
	return res
}

// Доломит
func getDolomite(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WDolo)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Доломит = ", sum)
	return sum
}

// Расход доломита, кг/т
func DolomiteFlow(db *sql.DB, date string) float64 {
	ccProd := prodMNLZ + prodIngot
	dolomite := getDolomite(db, date)
	res := SafeDivision(dolomite, (ccProd * 1000))

	logger.Debug("Доломит на плавку = Доломит / (Производство КЦ * 1000)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f * 1000)", res, dolomite, ccProd))
	return res
}

// Алюминий
func getAluminum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WAlrz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Алюминий = ", sum)
	return sum
}

// Алюминий на разогрев
func AluminumPreheating(db *sql.DB, date string) float64 {
	ccProd := prodMNLZ + prodIngot
	alu := getAluminum(db, date)
	res := SafeDivision(alu, ccProd)

	logger.Debug("Алюминий на разогрев = Алюминий / Производство КЦ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, alu, ccProd))
	return res
}

// Смесь
func getMixture(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WSmesrz)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Смесь = ", sum)
	return sum
}

// Смесь на плавку
func MixMelting(db *sql.DB, date string) float64 {
	mix := getMixture(db, date)
	ccProd := prodMNLZ + prodIngot
	res := SafeDivision(mix, ccProd)

	logger.Debug("Смесь на плавку = Смесь / Производство КЦ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, mix, ccProd))
	return res
}

// Si КЦ всего
func getFeSiCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Si КЦ всего = ", sum)
	return sum
}

// Si КЦ
func FeSiConsumption(db *sql.DB, date string) float64 {
	fesisum := getFeSiCC(db, date)
	ccProd := prodMNLZ + prodIngot
	res := SafeDivision((fesisum * 1000), ccProd)

	logger.Debug("Si КЦ = Si КЦ всего * 1000 / Производство КЦ")
	logger.Debug(fmt.Sprintf("%f = %f * 1000 / %f", res, fesisum, ccProd))
	return res
}

// FeSi65 расчет
func getSiModel(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Fesi65)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("FeSi65 расчет = ", sum)
	return sum
}

// Si по модели
func FeSiModelConsumption(db *sql.DB, date string) float64 {
	fesi65sum := getSiModel(db, date)
	res := SafeDivision(fesi65sum, prodMNLZ)

	logger.Debug("Si по модели = FeSi65 расчет / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, fesi65sum, prodMNLZ))
	return res
}

// SiMn КЦ всего
func getFeSiMnCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("SiMn КЦ всего = ", sum)
	return sum
}

// SiMn КЦ
func SiMnConsumption(db *sql.DB, date string) float64 {
	fesimnsum := getFeSiMnCC(db, date)
	cCprod := prodMNLZ + prodIngot
	res := SafeDivision((fesimnsum * 1000), cCprod)

	logger.Debug("SiMn КЦ = SiMn КЦ всего * 1000 / Производство КЦ")
	logger.Debug(fmt.Sprintf("%f = %f * 1000 / %f", res, fesimnsum, cCprod))
	return res
}

// SiMn расчет
func getSiMnModel(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Simn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("SiMn расчет = ", sum)
	return sum
}

// SiMn по модели
func SiMnModelConsumption(db *sql.DB, date string) float64 {
	simnmodelsum := getSiMnModel(db, date)

	res := SafeDivision(simnmodelsum, prodMNLZ)
	logger.Debug("SiMn по модели = SiMn расчет / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, simnmodelsum, prodMNLZ))
	return res
}

// Mn КЦ всего
func getMnCC(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum1 := Sum(data)

	logger.Debug("FeMn78 = ", sum1)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn88)
	data = database.ExecuteQuery(db, q)
	sum2 := Sum(data)

	logger.Debug("FeMn88 = ", sum2)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WFemn95)
	data = database.ExecuteQuery(db, q)
	sum3 := Sum(data)

	logger.Debug("FeMn95 = ", sum3)

	res := sum1 + sum2 + sum3

	logger.Debug("Mn КЦ всего = ", res)
	return res
}

// Mn КЦ
func FeMnConsumption(db *sql.DB, date string) float64 {
	mnCC := getMnCC(db, date)

	cCprod := prodMNLZ + prodIngot
	res := SafeDivision((mnCC * 1000), cCprod)
	logger.Debug("Mn КЦ = Mn КЦ всего * 1000 / Производство КЦ")
	logger.Debug(fmt.Sprintf("%f = %f * 1000 / %f", res, mnCC, cCprod))
	return res
}

// Mn расчет
func getMnModel(db *sql.DB, date string) (float64, float64) {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn78)
	data := database.ExecuteQuery(db, q)
	sum1 := Sum(data)

	logger.Debug("FeMn78 Расчет = ", sum1)

	q = fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.Femn88)
	data = database.ExecuteQuery(db, q)
	sum2 := Sum(data)

	logger.Debug("FeMn88 Расчет = ", sum2)

	return sum1, sum2
}

// Mn по модели
func FeMnModelConsumption(db *sql.DB, date string) float64 {
	feMn78, femn88 := getMnModel(db, date)
	mnmodel := feMn78 + femn88
	res := SafeDivision(mnmodel, prodMNLZ)
	logger.Debug("Mn по модели = (FeMn78 Расчет + FeMn88 Расчет) / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = (%f + %f) / %f", res, feMn78, femn88, prodMNLZ))
	return res
}

// Признак отсечки Монокон
func getMonokonTruncationCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.FlHeatCutoffSlag)
	data := database.ExecuteQuery(db, q)
	len := Len(data)

	logger.Debug("Признак отсечки Монокон = ", len)
	return len
}

// Доля плавок с отсечкой шлака
func SlagTruncationRatio(db *sql.DB, date string) float64 {
	MonokonCount := getMonokonTruncationCount(db, date)
	res := SafeDivision(MonokonCount, (sMeltCount + cPourCount))

	logger.Debug("Отсечка шлака = Количество признаков отсечки Монокон / (Количество плавок OHPC + Количество плавок разливка)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f)\n", res, MonokonCount, sMeltCount, cPourCount))
	return res
}

// Признак со скачиванием шлака
func getSkimmingFlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetSlagTruncation, date)
	data := database.ExecuteQuery(db, q)
	len := Len(data)

	logger.Debug("Признак со скачиванием шлака = ", len)
	return len
}

// Доля плавок со скачиванием шлака
func SlagSkimmingRatio(db *sql.DB, date string) float64 {
	SlagSkimmingCount := getSkimmingFlagCount(db, date)
	res := SafeDivision(SlagSkimmingCount, (sMeltCount + cPourCount))

	logger.Debug("Скачивание шлака = Количество признаков скачки / (Количество плавок OHPC + Количество плавок разливка)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f)", res, SlagSkimmingCount, sMeltCount, cPourCount))
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
	logger.Debug("Длительность плавки =  конец плавки - начало плавки")
	logger.Debug(fmt.Sprintf("Средняя длительность плавки = %f", res))
	return res
}

// % Fe в шлаке
func FePercentageInSlag(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.ShlMgo)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	logger.Debug(fmt.Sprintf("Fe в шлаке = %f", avg))
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
	slagCount := getFePercentageInSlagCount(db, date)

	res := SafeDivision(slagCount, (sMeltCount+cPourCount)) * 100
	logger.Debug("Кол-во шлаков = MgO / (Количество плавок OHPC + Количество плавок разливка) * 100")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f) * 100", res, slagCount, sMeltCount, cPourCount))
	return res
}

// w_FeMn78,  w_heat_mnlz > 0
func getWFeMn78MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn78, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeMn88, w_heat_mnlz > 0
func getWFeMn88MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn88)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn88, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeMn95, w_heat_mnlz > 0
func getWFeMn95MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFemn95)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn95, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeSi, w_heat_mnlz > 0
func getWFeSiMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeSi, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_SiMn, w_heat_mnlz > 0
func getWSiMnMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_SiMn, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr, w_heat_mnlz > 0
func getWFeCrMNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr025, w_heat_mnlz > 0
func getWFeCr025MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr025)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr025, w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr800, w_heat_mnlz > 0
func getWFeCr800MNLZ(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForMNLZ, date, c.GlobalConfig.Measurings.WFecr800)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr800, w_heat_mnlz > 0 = ", sum)
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

	logger.Debug("Всего феросплавов для МНЛЗ = ", res)
	return res
}

// w_FeMn78, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn78Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn78)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn78, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeMn88, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn88Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn88)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn88, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeMn95, w_iron > 0 && w_heat_mnlz > 0
func getWFeMn95Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFemn95)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeMn95, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeSi, w_iron > 0 && w_heat_mnlz > 0
func getWFeSiIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFesi)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeSi, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_SiMn, w_iron > 0 && w_heat_mnlz > 0
func getWSiMnIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WSimn)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_SiMn, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr, w_iron > 0 && w_heat_mnlz > 0
func getWFeCrIngot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr025, w_iron > 0 && w_heat_mnlz > 0
func getWFeCr025Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr025)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr025, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
	return sum
}

// w_FeCr800, w_iron > 0 && w_heat_mnlz > 0
func getWFeCr800Ingot(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetFerroalloysForIngots, date, c.GlobalConfig.Measurings.WFecr800)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("w_FeCr800, w_iron > 0 && w_heat_mnlz > 0 = ", sum)
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

	logger.Debug("Всего феросплавов для Слитков = ", res)
	return res
}

// Выход годного КЦ общий
func GoodCCOutput(db *sql.DB, date string) float64 {
	prodCC := prodMNLZ + prodIngot

	CCIronConsumption := mnlzConsumption + bcConsumption

	ferroalloys := ferroalloysIngot + ferroalloysMNLZ

	res := SafeDivision(prodCC, (CCIronConsumption + mnlzScrapConsumption + bcScrapConsumption + ferroalloys/1000))
	logger.Debug("Выход годного КЦ = Производство КЦ / (Потребление чугуна КЦ + Потребление лома МНЛЗ + Потребление лома Слиток + всего феросплавов КЦ /1000)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f + %f + %f / 1000)", res, prodCC, CCIronConsumption, mnlzScrapConsumption, bcScrapConsumption, ferroalloys))
	return res
}

// Выход годного КЦ МНЛЗ
func GoodCCMNLZOutput(db *sql.DB, date string) float64 {
	res := SafeDivision(mnlzConsumption, (mnlzScrapConsumption + ferroalloysMNLZ/1000))
	logger.Debug("Выход годного КЦ МНЛЗ = Потребление чугуна МНЛЗ / (Потребление лома МНЛЗ + всего феросплавов для МНЛЗ / 1000)")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f / 1000)", res, mnlzConsumption, mnlzScrapConsumption, ferroalloysMNLZ))
	return res
}

// Выход годного КЦ Слиток
func GoodCCIngotOutput(db *sql.DB, date string) float64 {
	res := SafeDivision(prodIngot, (bcConsumption + bcScrapConsumption + ferroalloysIngot))
	logger.Debug("Выход годного Слиток = Производство слитки / (Потребление чугуна слитки + Потребление лома Слиток + всего феросплавов для слитка))")
	logger.Debug(fmt.Sprintf("%f = %f / (%f + %f + %f)", res, prodIngot, bcConsumption, bcScrapConsumption, ferroalloysIngot))
	return res
}

// Получить УПК за номером
func getUPKByNumber(db *sql.DB, date string, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetUpk, date, c.GlobalConfig.Measurings.TmTreatment, n)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug(fmt.Sprintf("УПК%d = %f", n, avg))
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
	logger.Debug("Время обработки = (УПК1 + УПК2 + УПК3) / 3")
	logger.Debug(fmt.Sprintf("Время обработки = %f", res))
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
		logger.Debug(fmt.Sprintf("Ср. Нагрев УПК%d = %f", n, res[n-1]))
		res[n-1] = res[n-1] * 24
	}

	calc := floats.Sum(res)
	calc = calc / 3
	logger.Debug(fmt.Sprintf("Время дуги = %f", calc))
	return calc
}

// Известь
func getCaO(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCaoLf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Известь = ", sum)
	return sum
}

// Известь высоко-кальц
func getCaOHighCalcium(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCao90Lf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Известь высоко-кальц = ", sum)
	return sum
}

// Потребление Извести, кг/т
func LimestoneConsumption(db *sql.DB, date string) float64 {
	//prodMNLZ := productionMNLZSum(db, date)
	CaO := getCaO(db, date)
	CaOhigh := getCaOHighCalcium(db, date)

	res := SafeDivision((CaO + CaOhigh), prodMNLZ)
	logger.Debug("Потребление Извести = (Известь + Известь высоко-кальц) / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = (%f + %f) / %f", res, CaO, CaOhigh, prodMNLZ))
	return res
}

// Шпат
func getCaF2(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WCaf2Lf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("Шпат = ", sum)
	return sum
}

// Потребление шпата
func FluorsparConsumption(db *sql.DB, date string) float64 {
	//prodMNLZ := productionMNLZSum(db, date)
	caF2 := getCaF2(db, date)

	res := SafeDivision(caF2, prodMNLZ)
	logger.Debug("Потребление шпата = Шпат / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, caF2, prodMNLZ))
	return res
}

// АРК
func getAPK(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetGreaterThanZeroData, date, c.GlobalConfig.Measurings.WApkLf)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	logger.Debug("АРК = ", sum)
	return sum
}

// Потребление АРК
func ArgonOxygenConsumption(db *sql.DB, date string) float64 {
	//prodMNLZ := productionMNLZSum(db, date)
	apk := getAPK(db, date)

	res := SafeDivision(apk, prodMNLZ)
	logger.Debug("Потребление АРК = АРК / Производство МНЛЗ")
	logger.Debug(fmt.Sprintf("%f = %f / %f", res, apk, prodMNLZ))
	return res
}

// Электричество
func getElectricity(db *sql.DB, date string, id int, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetUpk, date, id, n)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)

	return sum
}

// Кол-во открытий
func getHeatStartCount(db *sql.DB, date string, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetInletTemperatureOxidation, date, c.GlobalConfig.Measurings.HeatStart, n)
	data := database.ExecuteQuery(db, q)
	len := Len(data)

	logger.Debug(fmt.Sprintf("Кол-во открытий УПК%d = %f", n, len))
	return len
}

// Расход Электричества
func ElectricityConsumption(db *sql.DB, date string) float64 {
	idEnergyList := []int{
		c.GlobalConfig.Measurings.Energy1,
		c.GlobalConfig.Measurings.Energy2,
		c.GlobalConfig.Measurings.Energy3,
		c.GlobalConfig.Measurings.Energy4,
		c.GlobalConfig.Measurings.Energy5,
		c.GlobalConfig.Measurings.Energy6,
		c.GlobalConfig.Measurings.Energy7,
		c.GlobalConfig.Measurings.Energy8,
	}

	sumEnergy := 0.0
	meltingCount := 0.0
	for n := 1; n <= 3; n++ {
		for _, id := range idEnergyList {
			sumEnergy += getElectricity(db, date, id, n)
		}
		logger.Debug(fmt.Sprintf("Электричество УПК%d = %f", n, sumEnergy))
		meltingCount += getHeatStartCount(db, date, n)
	}

	res := sumEnergy / meltingCount / 143

	logger.Debug("Расход Электричества = Электричество / Кол-во открытий / 143")
	logger.Debug(fmt.Sprintf("%f = %f / %f / 143", res, sumEnergy, meltingCount))
	return res
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

		logger.Debug(fmt.Sprintf("УПК%d Замер1 Т = %f", i, avg))
		res = res + avg
	}
	res = res / 3

	logger.Debug("Температура по приходу = (УПК1 + УПК2 + УПК3) / 3")
	logger.Debug("Температура по приходу = ", res)
	return res
}

// Окисленность по приходу
func InletOxidation(db *sql.DB, date string) float64 {
	res := 0.0
	for i := 1; i < 4; i++ {
		q := fmt.Sprintf(c.GlobalConfig.Querries.GetInletTemperatureOxidation, date, c.GlobalConfig.Measurings.O21, i)
		data := database.ExecuteQuery(db, q)
		avg := Avg(data)

		logger.Debug(fmt.Sprintf("УПК%d Замер1 О2 = %f", i, avg))
		res = res + avg
	}
	res = res / 3

	logger.Debug("Окисленность по приходу = (УПК1 + УПК2 + УПК3) / 3")
	logger.Debug("Окисленность по приходу = ", res)
	return res
}

// Количество шлаков
func getSlagCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.SampleTime)
	data := database.ExecuteQuery(db, q)
	len := Len(data)

	logger.Debug("Количество шлаков = ", len)
	return len
}

// Количество плавок
func getMeltingCount(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.HeatStart)
	data := database.ExecuteQuery(db, q)
	len := Len(data)

	logger.Debug("Количество плавок = ", len)
	return len
}

// Анализ шлаков УПК
func UPKSlagAnalysis(db *sql.DB, date string) float64 {
	slagCount := getSlagCount(db, date)
	meltCount := getMeltingCount(db, date)

	res := SafeDivision(slagCount, meltCount)
	logger.Debug("Анализ шлаков УПК = количество шлаков / количество плавок")
	logger.Debug(fmt.Sprintf("%f = %f / %f\n", res, slagCount, meltCount))
	return res
}

// МНЛЗ Открытие
func getOpening(db *sql.DB, date string, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.DtBegin, i)
	data := database.ExecuteQuery(db, q)

	if Len(data) > 0 {
		logger.Debug("МНЛЗ Открытие = ", data[0].Value)
	} else {
		logger.Debug("МНЛЗ Открытие = 0")
	}

	return data
}

// МНЛЗ Закрытие
func getClosing(db *sql.DB, date string, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.DtEnd, i)
	data := database.ExecuteQuery(db, q)

	if Len(data) > 0 {
		logger.Debug("МНЛЗ Закрытие = ", data[int(Len(data)-1)].Value)
	} else {
		logger.Debug("МНЛЗ Закрытие = 0")
	}
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
	logger.Debug("Цикл разливки =  Закрытие - Открытие")
	logger.Debug("Цикл разливки = ", res)
	return res
}

// Ср. скорость ручья
func getAvgSpeed(db *sql.DB, date string, i int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.AvgSpeed, i)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug(fmt.Sprintf("Ср. скорость ручья МНЛЗ%d = %f", i, avg))
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

	logger.Debug("Скорость разливки = Сумма(Ср. скорость ручья МНЛЗ(1 - 3)) / 3")
	logger.Debug("Скорость разливки = ", res)
	return res
}

// Расход смазки
func getGreaseConsumption(db *sql.DB, date string, id int, i int) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, id, i)
	data := database.ExecuteQuery(db, q)
	return data
}

// Конец серии
func getEndSeries(db *sql.DB, date string, i int) []float64 {
	var res []float64
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.EndSeries, i)
	data := database.ExecuteQuery(db, q)
	for _, v := range data {
		val, err := strconv.ParseFloat(*v.Value, 64)
		if err == nil {
			res = append(res, val)
		}
	}

	logger.Debug(fmt.Sprintf("Конец серии МНЛЗ%d = %f", i, res))
	return res
}

// Серийность
func getSerelization(db *sql.DB, date string, i int) []float64 {
	var res []float64
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.Serialization, i)
	data := database.ExecuteQuery(db, q)

	for _, v := range data {
		val, err := strconv.ParseFloat(*v.Value, 64)
		if err == nil {
			res = append(res, val)
		}
	}

	logger.Debug(fmt.Sprintf("Серийность МНЛЗ%d = %f", i, res))
	return res
}

// Серийность стопорной разливки
func CastingStopperSerial(db *sql.DB, date string) float64 {
	count, res := 0.0, 0.0
	for n := 1; n <= 3; n++ {
		var valueArrays [][]*float64

		for _, id := range idOilList {
			data := getGreaseConsumption(db, date, id, n)

			values := ParseFloatValues(data)

			valueArrays = append(valueArrays, values)
		}

		averages := CalculateAverages(valueArrays)
		endSeries := getEndSeries(db, date, n)
		serelization := getSerelization(db, date, n)

		if len(averages) == len(endSeries) {
			for i, v := range averages {
				if endSeries[i] == 1 && *v == 0 {
					value := serelization[i]
					res += value
					count++
				}
			}
		}
	}

	res = SafeDivision(res, count)

	logger.Debug("Серийность стопорной разливки = ", res)
	return res
}

// Функция для расчета открытой серийности
func calculateOpenSerial(averages []*float64, endSeries []float64, serelization []float64) float64 {
	count, res := 0.0, 0.0
	if len(averages) == len(endSeries) {
		for i, v := range averages {
			if endSeries[i] == 1 && *v > 0 {
				value := serelization[i]
				res += value
				count++
			}
		}
	}

	res = SafeDivision(res, count)

	return res
}

// Серийность открытой разливки МНЛЗ
func MNLZOpenSerial(db *sql.DB, date string, n int) float64 {
	var valueArrays [][]*float64

	for _, id := range idOilList {
		data := getGreaseConsumption(db, date, id, n)

		values := ParseFloatValues(data)

		valueArrays = append(valueArrays, values)
	}

	averages := CalculateAverages(valueArrays)
	endSeries := getEndSeries(db, date, n)
	serelization := getSerelization(db, date, n)

	res := calculateOpenSerial(averages, endSeries, serelization)

	logger.Debug(fmt.Sprintf("Серийность открытой разливки МНЛЗ%d = %f", n, res))
	return res
}

// Кол-во ручьев
func getStreamsCount(db *sql.DB, date string, n int) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.PCount, n)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)

	logger.Debug(fmt.Sprintf("Кол-во ручьев МНЛЗ%d = %f", n, avg))
	return avg
}

// Количество ручьев в работе МНЛЗ
func MNLZStreams(db *sql.DB, date string, n int) float64 {
	res := getStreamsCount(db, date, n)
	return res
}

// Перепаковка
func getRepackingMin(db *sql.DB, date string, n int) int {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.TmBetween, n)
	data := database.ExecuteQuery(db, q)
	seconds := Avg(data)
	minutes := int(math.Floor(seconds / 60))

	logger.Debug(fmt.Sprintf("Перепаковка МНЛЗ%d = %d", n, minutes))
	return minutes
}

// Длительность перепаковки МНЛЗ, мин
func MNLZRepackingDuration(db *sql.DB, date string, n int) float64 {
	res := getRepackingMin(db, date, n)
	return float64(res)
}

// Марка стали
func getSteelGrades(db *sql.DB, date string, n int) []*string {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.SteelGrade, n)
	data := database.ExecuteQuery(db, q)
	var grades []*string
	for _, v := range data {
		grades = append(grades, v.Value)
	}

	gradeStrings := make([]string, len(grades))
	for i, gradePtr := range grades {
		if gradePtr != nil {
			gradeStrings[i] = *gradePtr
		}
	}

	logger.Debug(fmt.Sprintf("Марки стали МНЛЗ%d = %s", n, strings.Join(gradeStrings, ", ")))
	return grades
}

// Температура МНЛЗ
func getTemperatureMnlz(db *sql.DB, date string, n int) float64 {
	var withinCount int
	var valueArrays [][]*float64
	idTempList := []int{
		c.GlobalConfig.Measurings.Temperature1,
		c.GlobalConfig.Measurings.Temperature2,
		c.GlobalConfig.Measurings.Temperature3,
		c.GlobalConfig.Measurings.Temperature4,
		c.GlobalConfig.Measurings.Temperature5,
		c.GlobalConfig.Measurings.Temperature6,
	}

	for _, t := range idTempList {
		q := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, t, n)
		data := database.ExecuteQuery(db, q)
		values := ParseFloatValues(data)

		valueArrays = append(valueArrays, values)
	}

	averages := CalculateAverages(valueArrays)
	steelGrades := getSteelGrades(db, date, n)
	serialization := getSerelization(db, date, n)

	for i, steel := range steelGrades {
		min, max := FindSteelGrade(*steel)
		logger.Debug(fmt.Sprintf("Марка стали = %s, Min = %d, Max = %d", *steel, min, max))
		logger.Debug(fmt.Sprintf("Серийность = %f, Температура = %f", serialization[i], *averages[i]))
		if serialization[i] == 1.0 || averages[i] == nil || (*averages[i] < float64(max) && *averages[i] > float64(min)) {
			withinCount++
		}
	}

	meltCount := getMeltingCount(db, date)
	res := 1.0 - SafeDivision(float64(withinCount), meltCount)

	logger.Debug(fmt.Sprintf("Плавки с отклонением по температуре МНЛЗ%d = кол-во норм / кол-во плавок", n))
	logger.Debug(fmt.Sprintf("%f = 1 - %d / %f", res, withinCount, meltCount))
	return res
}

// Плавки с отклонением по температуре МНЛЗ, %
func MNLZMeltTempDeviation(db *sql.DB, date string, n int) float64 {
	res := getTemperatureMnlz(db, date, n)
	return res
}

// Вес стали по ППС
func getWeightPPS(db *sql.DB, date string) float64 {
	wend, wstart := 0.0, 0.0
	for i := 1; i <= 3; i++ {

		q1 := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.WeightGrossEnd, i)
		data1 := database.ExecuteQuery(db, q1)

		q2 := fmt.Sprintf(c.GlobalConfig.Querries.GetMnlz, date, c.GlobalConfig.Measurings.WeightGrossStart, i)
		data2 := database.ExecuteQuery(db, q2)

		wend += Sum(data1)
		wstart += Sum(data2)
	}
	res := wstart - wend
	logger.Debug("Вес стали по ППС = ", res)
	return res
}

// Вес заготовки по SAP
func getWeightSAP(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.WeightLs)
	data := database.ExecuteQuery(db, q)
	res := Sum(data)

	logger.Debug("Вес заготовки по SAP = ", res)
	return res
}

// Вес заготовки по SAP
func GoodMNLZOutput(db *sql.DB, date string) float64 {
	wSap := getWeightSAP(db, date)
	wPps := getWeightPPS(db, date)

	res := 1 - SafeDivision(wSap, wPps)

	logger.Debug("Вес заготовки по SAP = 1 - Вес заготовки по SAP / Вес стали по ППС")
	logger.Debug(fmt.Sprintf("%f = 1 - %f / %f", res, wSap, wPps))
	return res
}

// Начало разливки МНЛЗ
func getStartTimeOfMNLZPouring(db *sql.DB, date string) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.CCMBegin)
	data := database.ExecuteQuery(db, q)

	if Len(data) > 0 {
		logger.Debug("Начало разливки МНЛЗ = ", data[0].Value)
	} else {
		logger.Debug("Начало разливки МНЛЗ = 0")
	}
	return data
}

// Время окончания выпуска
func getEndTimeOfProduction(db *sql.DB, date string) []models.Query {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetData, date, c.GlobalConfig.Measurings.DkTap)
	data := database.ExecuteQuery(db, q)

	if Len(data) > 0 {
		logger.Debug("Время окончания выпуска = ", data[0].Value)
	} else {
		logger.Debug("Время окончания выпуска = 0")
	}
	return data
}

// Время нахождения меиалла в ковше (до разливки), мин
func MetalRetentionTime(db *sql.DB, date string) float64 {
	dtn := getStartTimeOfMNLZPouring(db, date)
	dtk := getEndTimeOfProduction(db, date)
	if Len(dtn) == 0 || Len(dtk) == 0 {
		logger.Debug("Время нахождения металла в ковше (до разливки) = 0.0")
		return 0.0
	}
	res := AvgDiffDate(dtn, dtk)

	logger.Debug("Время нахождения металла в ковше (до разливки) = ", res)
	return res
}
