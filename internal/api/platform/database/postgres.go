package database

import (
	"fmt"
	"log"
	"os"

	"boardcamp/api/handler"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {

	var dsn string

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	sslrootcert := os.Getenv("DB_SSLROOTCERT")
	sslkey := os.Getenv("DB_SSLCLIENTKEY")
	sslcert := os.Getenv("DB_SSLCLIENTCERT")
	debug := os.Getenv("DB_DEBUG") == "true"

	// This is a function that is responsible for build database connecting string.
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/Sao_Paulo", host, user, pass, name, port)

	if sslmode != "" && sslrootcert != "" {
		dsn = fmt.Sprintf("%s sslmode=%s  sslrootcert=%s", dsn, sslmode, sslrootcert)
	}

	if sslkey != "" && sslcert != "" {
		dsn = fmt.Sprintf("%s sslkey=%s sslcert=%s", dsn, sslkey, sslcert)
	}

	log.Printf("Connecting to database: %s:%s\n", host, port)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Silent),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatalln(err)
	}

	if debug {
		DB = DB.Debug()
	}

	dbi, _ := DB.DB()
	err = dbi.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	DB.AutoMigrate(&handler.Game{}, &handler.Customer{}, &handler.Rental{})

	return DB
}
