package calculations

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/database"
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
	q := fmt.Sprintf(c.GlobalConfig.Querries.MeltdownsCasting, date)
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

	ICCasting := (MNLZConsumption + BCConsumption) / (SMeltCount + CPourCount)
	return ICCasting
}

// Расход лома на плавку
func ConsumptionOfScrapForMelting(db *sql.DB, date string) float64 {
	MNLZScrapConsumption := scrapConsumptionMnlzSum(db, date)
	BCScrapConsumption := scrapConsumptionIngotSum(db, date)
	SMeltCount := numberMeltdownsOnrs(db, date)
	CPourCount := numberMeltdownsCasting(db, date)

	ScrapMelting := (MNLZScrapConsumption + BCScrapConsumption) / (SMeltCount + CPourCount)
	return ScrapMelting
}

// Содержание Si в чугуне
func GetSiInCastIron(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.SiInCastIron, date)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Температура чугуна
func GetCastIronTemperature(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.CastIronTemperature, date)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Содержание S, %
func GetSContent(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.SContentPercentage, date)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Производство МНЛЗ
func productionMNLZSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.MnlzProduction, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Производство слитки
func productionIngotSum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.IngotProduction, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Средний вес плавки МНЛЗ, тонн
func MNLZMeltingAvgWeight(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	sMeltCount := numberMeltdownsOnrs(db, date)

	res := prodMNLZ / sMeltCount
	return res
}

// Средний вес плавки Слиток, тонн
func IngotMeltingAvgWeight(db *sql.DB, date string) float64 {
	prodIngot := productionIngotSum(db, date)
	cPourCount := numberMeltdownsCasting(db, date)

	res := prodIngot / cPourCount
	return res
}

// Содержание О2 на выпуске, ppm
func O2Content(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.O2ContentAtTheOutlet, date)
	data := database.ExecuteQuery(db, q)
	avg := Avg(data)
	return avg
}

// Известь вр.
func getLime(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetLime, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Расход извести, кг/т
func LimeFlow(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	KCprod := prodMNLZ + prodIngot

	lime := getLime(db, date)
	res := lime / (KCprod * 1000)
	return res
}

// Доломит
func getDolomite(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetDolomite, date)
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
	res := dolomite / (KCprod * 1000)
	return res
}

// Алюминий
func getAluminum(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetAluminum, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Алюминий на разогрев
func AluminumPreheating(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	KCprod := prodMNLZ + prodIngot

	alu := getAluminum(db, date)
	res := alu / KCprod
	return res
}

// Смесь
func getMixture(db *sql.DB, date string) float64 {
	q := fmt.Sprintf(c.GlobalConfig.Querries.GetAluminum, date)
	data := database.ExecuteQuery(db, q)
	sum := Sum(data)
	return sum
}

// Смесь на плавку
func MixMelting(db *sql.DB, date string) float64 {
	prodMNLZ := productionMNLZSum(db, date)
	prodIngot := productionIngotSum(db, date)
	KCprod := prodMNLZ + prodIngot

	mix := getMixture(db, date)
	res := mix / KCprod
	return res
}
