package database

import (
	"database/sql"
	"fmt"

	config "github.com/dabi-ngin/discgo-bot/Config"
	logger "github.com/dabi-ngin/discgo-bot/Logger"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func Init() bool {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.ServiceSettings.DB_USER,
		config.ServiceSettings.DB_PASSWORD,
		config.ServiceSettings.DB_IP_ADDRESS,
		config.ServiceSettings.DB_PORT,
		config.ServiceSettings.DB_NAME)

	dbIn, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		logger.Error("", err)
		return false
	}

	Db = dbIn
	return true
}
