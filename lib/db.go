package lib

import (
	"fmt"
	"log"
	"os"
	"time"

	entity "github.com/LinggaAskaEdo/burger-duls/models/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(zapLogger Logger) Database {
	username := os.Getenv("MARIA_DB_USER")
	password := os.Getenv("MARIA_DB_PASS")
	host := os.Getenv("MARIA_DB_HOST")
	port := os.Getenv("MARIA_DB_PORT")
	dbname := os.Getenv("BURGER_DULS_DB_NAME")

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		zapLogger.Info("Url: ", url)
		zapLogger.Panic(err)
	}

	zapLogger.Info("Database connection established")

	db.AutoMigrate(&entity.User{}, &entity.Menu{}, &entity.Transaction{}, &entity.DetailTransaction{})

	zapLogger.Info("Table(s) migration was successful")

	return Database{
		DB: db,
	}
}
