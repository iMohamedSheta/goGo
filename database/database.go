package database

import (
	"database/sql"
	"fmt"
	"imohamedsheta/gocrud/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	config := config.GetDefaultDatabaseConfig()
	driver := config["driver"].(string)

	var dsn string

	switch driver {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
			config["user"],
			config["pass"],
			config["host"],
			config["port"],
			config["db_name"],
			config["charset"],
		)
	case "pgsql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config["host"],
			config["port"],
			config["user"],
			config["pass"],
			config["db_name"],
		)

	default:
		log.Fatalf("❌ Unsupported database driver: %s", driver)
	}

	var err error

	DB, err = sql.Open(driver, dsn)

	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %s", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database:  %s", err)
	}

	fmt.Println("✅ Connected to database")
}
