package ewc

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BaseDbService struct {
	ErrorsChan chan string
}

func (srv *BaseDbService) dbExec(closure func(db *gorm.DB)) {
	db := srv.getDb()

	if db == nil {
		return
	}

	closure(db)

	if err := db.Close(); err != nil {
		log.Println("close db error:", err)
	}
}

func (srv *BaseDbService) getDb() *gorm.DB {
	/*if driver == "sqlite3" {
		if _, err := os.Stat(connectionString); os.IsNotExist(err) {
			file, err := os.Create(connectionString)

			if err != nil {
				log.Println("create db file error:", err)
				log.Println(connectionString)
				return nil
			}
			if err := file.Close(); err != nil {
				log.Println("close db file error:", err)
			}
		}
	}*/

	db, err := gorm.Open(driver, connectionString)

	if err != nil {
		log.Println("open db error:", err.Error())
		return nil
	}

	return db
}
