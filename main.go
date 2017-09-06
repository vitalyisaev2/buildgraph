package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/zap"

	"github.com/vitalyisaev2/buildgraph/db/models"
)

func createDatabase(logger *zap.SugaredLogger) (*gorm.DB, error) {

	// init database
	logger.Debug("Open database file")
	db, err := gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		return nil, err
	}

	// register models
	logger.Debug("Register models")
	models := []interface{}{
		&models.Project{},
		&models.Build{},
		&models.Author{},
		&models.Hook{},
		&models.Commit{},
		&models.PathEvent{},
		&models.Path{},
		&models.PathChange{},
	}
	if err := db.AutoMigrate(models...).Error; err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// init logger
	logger := zap.NewExample().Sugar()
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	db, err := createDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()
}
