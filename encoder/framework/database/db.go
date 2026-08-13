package database

import (
	"enconder/domain"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Db            *gorm.DB
	Dsn           string
	DsnTest       string
	DbType        string
	DbTypeTest    string
	Debug         bool
	AutoMigrateDb bool
	Env           string
}

func NewDb() *Database {
	return &Database{}
}

func NewDbTest() *gorm.DB {
	dbInstance := NewDb()
	dbInstance.Env = "test"
	dbInstance.DbTypeTest = "sqlite3"
	dbInstance.DsnTest = ":memory:"
	dbInstance.AutoMigrateDb = true
	dbInstance.Debug = true

	connection, err := dbInstance.Connect()
	if err != nil {
		log.Fatalf("Test db error: %v", err)
	}

	return connection
}

func (d *Database) Connect() (*gorm.DB, error) {
	var err error
	var dialector gorm.Dialector

	if d.Env == "test" {
		dialector = sqlite.Open(d.DsnTest)
	} else {
		if d.Dsn == "" {
			d.Dsn = "gorm.db"
		}
		dialector = sqlite.Open(d.Dsn)
	}

	d.Db, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar no SQLite: %v", err)
	}

	if d.Debug {
		d.Db.Logger = d.Db.Logger.LogMode(logger.Info)
	} else {
		d.Db.Logger = d.Db.Logger.LogMode(logger.Silent)
	}

	if d.AutoMigrateDb {
		err = d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		if err != nil {
			return nil, fmt.Errorf("falha ao rodar automigrate: %v", err)
		}
	}

	return d.Db, nil
}

