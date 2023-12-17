package database

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sayantan-s/easy-send/config"
	"github.com/sayantan-s/easy-send/models"
)

var (
    db     *gorm.DB
    dbOnce sync.Once
)

func GetInstance() (*gorm.DB, error) {
    var err error

	dbOnce.Do(func() {
        db, err = initDB()
    })

    return db, err
}

func initDB() (*gorm.DB, error) {
	connectionString := config.GetConfig("PG_CONNECTION_STRING")
    
	db, err := gorm.Open("postgres", connectionString)
    
	if err != nil {
        return nil, err
    }

    models.AutoMigrateTranscripts(db)

    return db, nil
}