package database

import (
	"context"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"gorm.io/gorm"
)

type CampusDb struct {
	db *gorm.DB
}

func NewCampusDb(db *gorm.DB) *CampusDb {
	return &CampusDb{db: db}
}

func (c CampusDb) CreateMany(ctx context.Context, campuses []domain.Campus) error {
	return c.db.WithContext(ctx).Create(&campuses).Error
}

func (c CampusDb) GetByID(ctx context.Context, campusID uint) (domain.Campus, error) {
	var campus domain.Campus
	err := c.db.WithContext(ctx).First(&campus, campusID).Error
	if err == gorm.ErrRecordNotFound {
		return domain.Campus{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Campus{}, err
	}
	return campus, nil
}

func (c CampusDb) GetByName(ctx context.Context, collegeId uint, name string) ([]domain.Campus, error) {
	var campuses []domain.Campus
	err := c.db.WithContext(ctx).Where(
		domain.Campus{CollegeID: collegeId, Name: name},
	).Find(&campuses).Error
	if err != nil {
		return nil, err
	}
	return campuses, nil
}

func (c CampusDb) GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.Campus, error) {
	var campuses []domain.Campus
	if err := c.db.WithContext(ctx).
		Where("college_id = ?", collegeID).
		Find(&campuses).Error; err != nil {
		return nil, err
	}
	return campuses, nil

}

func (c CampusDb) GetByCollegeIDs(ctx context.Context, collegeIDs []uint) ([]domain.Campus, error) {
	var campuses []domain.Campus
	if err := c.db.WithContext(ctx).
		Where("college_id IN ?", collegeIDs).
		Find(&campuses).Error; err != nil {
		return nil, err
	}
	return campuses, nil

}

func (c CampusDb) WithTx(tx repository.Tx) repository.CampusRepo {
	return CampusDb{db: tx.DB().(*gorm.DB)}
}
