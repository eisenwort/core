package ewc

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type BaseDbService struct {
	driver           string
	connectionString string
	ErrorsChan       chan string
}

func (srv *BaseDbService) dbExec(closure func(db *gorm.DB)) {
	db := srv.getDb()

	if db == nil {
		return
	}

	closure(db)

	db.Close()
}

func (srv *BaseDbService) getDb() *gorm.DB {
	if srv.driver == "sqlite3" {
		if _, err := os.Stat(connectionString); os.IsNotExist(err) {
			file, err := os.Create(connectionString)

			if err != nil {
				log.Println("create db file error:", err)
				log.Println(connectionString)
				return nil
			}

			file.Close()
		}
	}

	db, err := gorm.Open(srv.driver, srv.connectionString)

	if err != nil {
		log.Println("open db error:", err.Error())
		return nil
	}

	db = db.Debug()
	return db
}
