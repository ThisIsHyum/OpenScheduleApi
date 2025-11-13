package database

import (
	"fmt"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Db struct {
	db *gorm.DB
}

func NewDb(config *config.Config) (*Db, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True",
		config.Db.User, config.Db.Password,
		config.Db.Host, config.Db.Port, config.Db.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Parser{}, &models.College{}, &models.Campus{},
		&models.Lesson{}, &models.StudentGroup{}, &models.Call{},
	); err != nil {
		return nil, fmt.Errorf("auto migration failed: %w", err)
	}
	return &Db{db: db}, nil
}
