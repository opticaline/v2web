package core

import (
	"fmt"
	"log"
)
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Database struct {
	*gorm.DB
}

func Create() *Database {
	db := &Database{Init()}
	db.createTables()
	return db
}

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)
	return db
}

func (db *Database) createTables() {
	log.Println("Check Tables")
	if !db.HasTable(&Traffic{}) {
		db.CreateTable(&Traffic{})
	}
}
