package ewc

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func dbExec(closure func(db *gorm.DB)) {
	if currentUser.Reseted {
		return
	}
	db := getDb()

	if db == nil {
		return
	}

	closure(db)
}

func getDb() *gorm.DB {
	if db != nil {
		return db
	}

	db, err := gorm.Open(driver, connectionString)

	if err != nil {
		log.Println("open db error:", err.Error())
		return nil
	}
	if driver != "sqlite3" {
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		db.DB().SetConnMaxLifetime(time.Hour)
	}

	return db
}
