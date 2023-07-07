package database

import (
	"database/sql"
	"fmt"
	"log"
	c "main/configuration"
	"main/logger"
	"main/models"

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
		var data models.Query
		err := rows.Scan(&data.IdMeasuring, &data.TimeStamp, &data.Value, &data.Quality, &data.BatchId)
		if err != nil {
			fmt.Println("Failed to scan row:", err)
			continue
		}
		query = append(query, data)
	}
	fmt.Println(query)
	return query
}
