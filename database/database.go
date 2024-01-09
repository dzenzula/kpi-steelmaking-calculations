package database

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/logger"
	"main/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Masterminds/structable"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

var retry string = "Next try to connect wil be in 5min"

func ConnectToDatabase(config models.ConStringPG, dbName string) *sql.DB {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, dbName, config.SSLMode,
	)

	for {
		conn, conErr := sql.Open(c.GlobalConfig.TypePG, connString)
		if conErr != nil {
			logger.Error("Error opening database connection:", conErr.Error())
			logger.Error(retry)
			time.Sleep(5 * time.Minute)
			continue
		}

		pingErr := conn.Ping()
		if pingErr != nil {
			logger.Error("Error pinging database:", pingErr.Error())
			logger.Error(retry)
			time.Sleep(5 * time.Minute)
			continue
		}

		logger.Info(fmt.Sprintf("Connected to %s, %s", config.Host, dbName))
		return conn
	}
}

func ConnectPg() *sql.DB {
	return ConnectToDatabase(c.GlobalConfig.ConStringPgDb, c.GlobalConfig.ConStringPgDb.DBName)
}

func ExecuteQuery(db *sql.DB, q string) []models.Query {
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var query []models.Query
	for rows.Next() {
		var valueNullable sql.NullString
		var data models.Query
		err := rows.Scan(&data.IdMeasuring, &data.TimeStamp, &valueNullable, &data.Quality, &data.BatchId, &data.Timestamp_insert)
		if err != nil {
			logger.Error("Failed to scan row:", err)
			continue
		}

		if valueNullable.Valid {
			data.Value = &valueNullable.String
		} else {
			// Если значение null, устанавливаем значение Value в nil
			data.Value = nil
		}

		query = append(query, data)
	}
	return query
}

func InsertMsReport(db *sql.DB, report models.Report) {
	runner := squirrel.NewStmtCacheProxy(db)

	// Создание объекта structable для привязки значений структуры
	//mydb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).RunWith(runner)

	r := structable.New(runner, c.GlobalConfig.TypeMS).Bind("[dbo].[KpiReport]", report)
	err := r.Insert()
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Data inserted in report MsSQL!")
}

func InsertPgReport(db *sql.DB, report models.Report) {
	runner := squirrel.NewStmtCacheProxy(db)
	r := structable.New(runner, c.GlobalConfig.TypePG).Bind("reports.\"kpi-steelmaking-reports\"", &report)
	err := r.Insert()
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Data inserted in report PostgreSQL!")
}

func UpdateMsReport(db *sql.DB, report models.Report) {
	runner := squirrel.NewStmtCacheProxy(db)
	r := structable.New(runner, c.GlobalConfig.TypePG).Bind("[dbo].[KpiReport]", &report)
	err := r.Update()
	if err != nil {
		logger.Error(err)
		r.Insert()
		logger.Info("Current year data inserted in MsSQL!")
		return
	}

	logger.Info("Current year data updated in MsSQL!")
}

func UpdatePgReport(db *sql.DB, report models.Report) int {
	runner := squirrel.NewStmtCacheProxy(db)
	r := structable.New(runner, c.GlobalConfig.TypePG).Bind("reports.\"kpi-steelmaking-reports\"", &report)

	if report.Id == 0 {
		err := r.Insert()
		if err != nil {
			logger.Error(err)
		}
		logger.Info("Current year data inserted in PostgreSQL!")
	} else {
		err := r.Update()
		if err != nil {
			logger.Error(err)
		}
		logger.Info("Current year data updated in PostgreSQL!")
	}

	return report.Id
}
