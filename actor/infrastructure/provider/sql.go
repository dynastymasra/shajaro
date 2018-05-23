package provider

import (
	"sirius/actor/config"

	"github.com/jinzhu/gorm"
	"github.com/matryer/resync"

	"reflect"
	"runtime"

	"fmt"

	"errors"

	log "github.com/dynastymasra/gochill"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

var (
	db      *gorm.DB
	err     error
	runOnce resync.Once
)

func ConnectSQL() (*gorm.DB, error) {
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.DatabaseUsername,
		config.DatabasePassword, config.DatabaseAddress, config.DatabaseName)

	runOnce.Do(func() {
		db, err = gorm.Open("postgres", databaseURL)

		if err != nil {
			log.Error(log.Msg("Failed connect to database", err.Error()), log.O("version", config.Version),
				log.O("package", runtime.FuncForPC(reflect.ValueOf(ConnectSQL).Pointer()).Name()),
				log.O("project", config.ProjectName), log.O("database_url", databaseURL))
			return
		}

		db.DB().SetMaxIdleConns(config.DatabaseMaxIdle)
		db.DB().SetMaxOpenConns(config.DatabaseMaxOpen)

		err = db.DB().Ping()
		if err != nil {
			log.Error(log.Msg("Failed ping to database", err.Error()), log.O("version", config.Version),
				log.O("package", runtime.FuncForPC(reflect.ValueOf(ConnectSQL).Pointer()).Name()),
				log.O("project", config.ProjectName), log.O("database_url", databaseURL))
			return
		}

		db.LogMode(config.DatabaseLog)
	})

	return db, err
}

func Migration(data *gorm.DB) error {
	driver, err := postgres.WithInstance(data.DB(), &postgres.Config{})
	if err != nil {
		log.Error(log.Msg("Failed open instance", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(Migration).Pointer()).Name()),
			log.O("project", config.ProjectName))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migration", "postgres", driver)
	if err != nil {
		log.Error(log.Msg("Failed open database", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(Migration).Pointer()).Name()),
			log.O("project", config.ProjectName))
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(log.Msg("Failed migrate database", err.Error()), log.O("version", config.Version),
			log.O("package", runtime.FuncForPC(reflect.ValueOf(Migration).Pointer()).Name()),
			log.O("project", config.ProjectName))
		m.Down()
		return err
	}

	return nil
}

func SQLPing(db *gorm.DB) error {
	if db == nil {
		return errors.New(config.ErrDatabaseNil)
	}
	return db.DB().Ping()
}

func CloseDB(db *gorm.DB) error {
	if db == nil {
		return errors.New(config.ErrDatabaseNil)
	}
	return db.Close()
}

func ResetDBSingleton() {
	runOnce.Reset()
}
