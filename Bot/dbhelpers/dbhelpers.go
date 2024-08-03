package dbhelper

import (
	"database/sql"
	"fmt"

	"github.com/dabi-ngin/discgo-bot/Bot/audit"
	"github.com/dabi-ngin/discgo-bot/Bot/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init: runs on package declaration when another module brings db module into scope.
// does not need to be called directly.
// sets up the MySQL session connection and logs details of success/failure.
func LoadDB() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_IP_ADDRESS,
		config.DB_PORT,
		config.DB_NAME)
	audit.Log("Opening DB Connection: " + dataSourceName)
	dbIn, err := sql.Open("mysql", dataSourceName)
	// dbIn, err := sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbIP+":"+dbPort+")/"+dbName+"?parseTime=true")
	if err != nil {
		audit.Error(err)
	} else {
		audit.Log("DB Connection successful")
		db = dbIn
	}
}

func GetDB() *sql.DB {
	if db == nil {
		LoadDB()
	}
	return db
}
