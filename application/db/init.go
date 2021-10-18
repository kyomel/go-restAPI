package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// Postgres ...
func DBInit() *gorm.DB {
	log.Info("Starting Database Connection")
	dbURI := "host=localhost user=kyomel password=limabelas15 dbname=todo_list_go port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Panic("failed to connect database with error " + err.Error())
	}
	db.LogMode(true)
	return db
}
