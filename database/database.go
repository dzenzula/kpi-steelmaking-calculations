package database

import (
	"database/sql"
	"fmt"
	c "main/configuration"
	"main/logger"
	"main/models"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/Masterminds/structable"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

var (
	connMutex sync.Mutex // Мьютекс для обеспечения безопасности доступа к подключениям
)

func ConnectMs() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", c.GlobalConfig.ConStringMsDb.Server, c.GlobalConfig.ConStringMsDb.UserID, c.GlobalConfig.ConStringMsDb.Password, c.GlobalConfig.ConStringMsDb.Database)

	connMutex.Lock()
	defer connMutex.Unlock()

	conn, conErr := sql.Open(c.GlobalConfig.TypeMS, connString)
	if conErr != nil {
		logger.Error("Error opening database connection:", conErr.Error())
		return nil
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		logger.Fatal(pingErr.Error())
	}

	logger.Info("Connected to", c.GlobalConfig.ConStringMsDb.Server, c.GlobalConfig.ConStringMsDb.Database)
	return conn
}

func ConnectToDatabase(config models.ConStringPG, dbName string) *sql.DB {
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, dbName, config.SSLMode,
	)

	connMutex.Lock()
	defer connMutex.Unlock()

	conn, conErr := sql.Open(c.GlobalConfig.TypePG, connString)
	if conErr != nil {
		return nil
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		return nil
	}

	logger.Info(fmt.Sprintf("Connected to %s, %s\n", config.Host, dbName))
	return conn
}

func ConnectPgData() *sql.DB {
	return ConnectToDatabase(c.GlobalConfig.ConStringPgDb, c.GlobalConfig.ConStringPgDb.DBName)
}

func ConnectPgReports() *sql.DB {
	return ConnectToDatabase(c.GlobalConfig.ConStringPgReports, c.GlobalConfig.ConStringPgReports.DBName)
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
			fmt.Println("Failed to scan row:", err)
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

	logger.Info("Data to report inserted!")
}

func InsertPgReport(db *sql.DB, report models.Report) {
	runner := squirrel.NewStmtCacheProxy(db)
	r := structable.New(runner, c.GlobalConfig.TypePG).Bind("reports.\"kpi-steelmaking-reports\"", &report)
	err := r.Insert()
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Data to report inserted in PostgreSQL!")
}
