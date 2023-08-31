package database

import (
	"database/sql"
	"fmt"
	"log"
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
	msDevConn *sql.DB    // Глобальная переменная для подключения к MS SQL
	pgDevConn *sql.DB    // Глобальная переменная для подключения к PostgreSQL
	connMutex sync.Mutex // Мьютекс для обеспечения безопасности доступа к подключениям
)

func ConnectMsDev() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", c.GlobalConfig.ConStringMsDev.Server, c.GlobalConfig.ConStringMsDev.UserID, c.GlobalConfig.ConStringMsDev.Password, c.GlobalConfig.ConStringMsDev.Database)

	connMutex.Lock()
	defer connMutex.Unlock()

	conn, conErr := sql.Open(c.GlobalConfig.TypeMS, connString)
	if conErr != nil {
		logger.Error("Error opening database connection:", conErr.Error())
		return nil
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr.Error())
	}

	msDevConn = conn

	return msDevConn
}

func ConnectPgDev() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.GlobalConfig.ConStringPgDev.Host, c.GlobalConfig.ConStringPgDev.Port, c.GlobalConfig.ConStringPgDev.User, c.GlobalConfig.ConStringPgDev.Password, c.GlobalConfig.ConStringPgDev.DBName, c.GlobalConfig.ConStringPgDev.SSLMode)

	connMutex.Lock()
	defer connMutex.Unlock()

	conn, conErr := sql.Open(c.GlobalConfig.TypePG, connString)
	if conErr != nil {
		logger.Error("Error opening database connection:", conErr.Error())
		return nil
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		logger.Fatal(pingErr.Error())
	}

	pgDevConn = conn

	return pgDevConn
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
		err := rows.Scan(&data.IdMeasuring, &data.TimeStamp, &valueNullable, &data.Quality, &data.BatchId)
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

func InsertReport(db *sql.DB, report models.Report) error {
	runner := squirrel.NewStmtCacheProxy(db)

	// Создание объекта structable для привязки значений структуры
	//mydb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question).RunWith(runner)

	r := structable.New(runner, c.GlobalConfig.TypeMS).Bind("[dbo].[KpiReport]", report)
	err := r.Insert()
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
