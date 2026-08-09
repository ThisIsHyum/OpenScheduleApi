package database

import (
	"fmt"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Tx struct {
	db *gorm.DB
}

func NewDb(cfgDb *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True",
		cfgDb.User, cfgDb.Password,
		cfgDb.Host, cfgDb.Port, cfgDb.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func InitCreateTx(db *gorm.DB) repository.CreateTx {
	return func() (repository.Tx, error) {
		db := db.Begin()
		return Tx{db: db}, db.Error
	}
}

func (t Tx) Commit() error   { return t.db.Commit().Error }
func (t Tx) Rollback() error { return t.db.Rollback().Error }
func (t Tx) DB() any         { return t.db }

func NewRepos(db *gorm.DB) repository.Repos {
	return repository.Repos{
		CollegeRepo:      NewCollegeDb(db),
		CallRepo:         NewCallDb(db),
		StudentGroupRepo: NewGroupDb(db),
		LessonRepo:       NewLessonDb(db),
		CampusRepo:       NewCampusDb(db),
		CreateTx:         InitCreateTx(db),
	}
}
