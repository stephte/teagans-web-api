package config

import (
	"teagans-web-api/app/models"
	"gorm.io/driver/postgres"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"fmt"
)

// -------- DB Setup --------

type DBConn struct {
	host			string
	user			string
	password		string
	name			string
	port			int

	db				*gorm.DB
	logger			zerolog.Logger

	verbose			bool
}

func InitDBConn(logger zerolog.Logger, verbose bool) DBConn {
	return DBConn{
		logger: logger,
		verbose: verbose,
	}
}


func(this *DBConn) FireUp() (error) {
	if this.verbose {
		this.logger.Info().Msg("Firing Up DB!")
	}

	connstring := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		this.host, this.user, this.password, this.name, this.port,
	)

	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})
	if err != nil {
		return err
	}

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	this.db = db
	this.migrate()

	if this.verbose {
		this.logger.Info().Msg("DB fired up!")
	}

	return nil
}


func(this DBConn) CoolDown() error {
	if this.verbose {
		this.logger.Info().Msg("Cooling down DB")
	}

	sqlDB, sqlErr := this.db.DB()
	
	if sqlErr != nil {
		// not much else to do if the app is already shutting down and theres an error here...
		this.logger.Error().Err(sqlErr).Msg("Error getting DB connection")
		return sqlErr
	}

	closeErr := sqlDB.Close()
	if closeErr != nil {
		this.logger.Error().Err(closeErr).Msg("Error closing DB connection")
		return closeErr
	}

	return nil
}


func(this *DBConn) migrate() {

	// this.logger.Warn().Msg("Dropping Users table")
	// this.db.Migrator().DropTable(&models.User{})

	if this.verbose {
		this.logger.Info().Msg("Migrating...")
	}
	// add DB table models here
	this.db.AutoMigrate(
		&models.User{},
	)
}

// ---- setters ----

func(this *DBConn) SetHost(host string) {
	this.host = host
}

func(this *DBConn) SetUser(user string) {
	this.user = user
}

func(this *DBConn) SetPassword(password string) {
	this.password = password
}

func(this *DBConn) SetName(name string) {
	this.name = name
}

func(this *DBConn) SetPort(port int) {
	this.port = port
}

// ---- Getters ----

func(this DBConn) GetDB() *gorm.DB {
	return this.db
}
