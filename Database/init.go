package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	config "github.com/hashbat-dev/discgo-bot/Config"
	logger "github.com/hashbat-dev/discgo-bot/Logger"
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
	logger.Db = dbIn // Need a second instance in Logger to avoid import cycles, fun
	return true
}
