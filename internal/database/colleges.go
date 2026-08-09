package database

import (
	"context"
	"errors"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"

	"gorm.io/gorm"
)

type CollegeDb struct {
	db *gorm.DB
}

func NewCollegeDb(db *gorm.DB) *CollegeDb {
	return &CollegeDb{db: db}
}

func (c CollegeDb) Create(ctx context.Context, college domain.College) (uint, error) {
	return college.ID, c.db.WithContext(ctx).Create(&college).Error
}

func (c CollegeDb) Get(ctx context.Context, collegeID uint) (domain.College, error) {
	var college domain.College
	err := c.db.WithContext(ctx).First(&college, collegeID).Error
	if err == gorm.ErrRecordNotFound {
		return domain.College{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.College{}, err
	}
	return college, nil
}

func (c CollegeDb) GetByName(ctx context.Context, name string) ([]domain.College, error) {
	var colleges []domain.College
	err := c.db.WithContext(ctx).Where("name = ?", name).Find(&colleges).Error
	if err != nil {
		return nil, err
	}
	return colleges, nil
}

func (c CollegeDb) GetAll(ctx context.Context) ([]domain.College, error) {
	var colleges []domain.College
	err := c.db.WithContext(ctx).Find(&colleges).Error
	if err != nil {
		return nil, err
	}
	return colleges, nil
}

func (c CollegeDb) GetByGroupID(ctx context.Context, groupID uint) (domain.College, error) {
	var college domain.College
	err := c.db.WithContext(ctx).Table("student_groups AS sg").
		Select("c.id, c.name").
		Joins("JOIN campuses AS c ON c.id = sg.campus_id").
		Joins("JOIN colleges AS co ON co.id = c.college_id").
		Where("sg.id = ?", groupID).
		Scan(&college).Error
	if err == gorm.ErrRecordNotFound {
		return domain.College{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.College{}, err
	}
	return college, nil
}

func (c CollegeDb) GetByToken(ctx context.Context, token string) (domain.College, error) {
	var college domain.College
	if err := c.db.WithContext(ctx).Where("token = ?", token).First(&college).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.College{}, domain.ErrNotFound
	} else if err != nil {
		return domain.College{}, err
	}
	return college, nil
}

func (c CollegeDb) Delete(ctx context.Context, collegeID uint) error {
	result := c.db.WithContext(ctx).Delete(&domain.College{}, collegeID)
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (c CollegeDb) WithTx(tx repository.Tx) repository.CollegeRepo {
	return CollegeDb{db: tx.DB().(*gorm.DB)}
}
