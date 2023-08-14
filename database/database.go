package database

import (
	"database/sql"
	"fmt"
	"log"
	c "main/configuration"
	"main/logger"
	"main/models"

	"github.com/Masterminds/squirrel"
	"github.com/Masterminds/structable"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

func ConnectMsDev() *sql.DB {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", c.GlobalConfig.ConStringMsDev.Server, c.GlobalConfig.ConStringMsDev.UserID, c.GlobalConfig.ConStringMsDev.Password, c.GlobalConfig.ConStringMsDev.Database)
	conn, conErr := sql.Open(c.GlobalConfig.TypeMS, connString)
	if conErr != nil {
		logger.Error("Error opening database connection:", conErr.Error())
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr.Error())
	}

	return conn
}

func ConnectPgDev() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.GlobalConfig.ConStringPgDev.Host, c.GlobalConfig.ConStringPgDev.Port, c.GlobalConfig.ConStringPgDev.User, c.GlobalConfig.ConStringPgDev.Password, c.GlobalConfig.ConStringPgDev.DBName, c.GlobalConfig.ConStringPgDev.SSLMode)
	conn, conErr := sql.Open(c.GlobalConfig.TypePG, connString)
	if conErr != nil {
		logger.Error("Error opening database connection:", conErr.Error())
	}

	pingErr := conn.Ping()
	if pingErr != nil {
		logger.Fatal(pingErr.Error())
	}

	return conn
}

func ExecuteQuery(db *sql.DB, q string) []models.Query {
	rows, err := db.Query(q)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	var query []models.Query
	for rows.Next() {
		var valueNullable sql.NullFloat64
		var data models.Query
		err := rows.Scan(&data.IdMeasuring, &data.TimeStamp, &valueNullable, &data.Quality, &data.BatchId)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}

		if valueNullable.Valid {
			data.Value = &valueNullable.Float64
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
