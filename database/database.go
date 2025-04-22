package database

import (
	"fmt"
	"imohamedsheta/gocrud/pkg/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	// _ "github.com/lib/pq"
)

var db *sqlx.DB

func Connect() {
	config := config.App.Get("database").(map[string]any)
	defaultDatabaseConnection := config["default"].(string)
	connectionConfig := config["connections"].(map[string]any)[defaultDatabaseConnection].(map[string]any)
	driver := connectionConfig["driver"].(string)

	var dsn string

	switch driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			connectionConfig["user"],
			connectionConfig["pass"],
			connectionConfig["host"],
			connectionConfig["port"],
			connectionConfig["database"],
			connectionConfig["charset"],
		)
	// case "pgsql":
	// 	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 		connectionConfig["host"],
	// 		connectionConfig["port"],
	// 		connectionConfig["user"],
	// 		connectionConfig["pass"],
	// 		connectionConfig["database"],
	// 		connectionConfig["sslmode"],
	// 	)

	default:
		log.Fatalf("❌ Unsupported database driver: %s", driver)
	}

	var err error

	db, err = sqlx.Open(driver, dsn)

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database:  %s", err)
	}

	fmt.Println("✅ Connected to database")
}

func DB() *sqlx.DB {
	return db
}
