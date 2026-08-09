package database

import (
	"context"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GroupDb struct {
	db *gorm.DB
}

func NewGroupDb(db *gorm.DB) *GroupDb {
	return &GroupDb{db: db}
}

func (g GroupDb) UpsertMany(ctx context.Context, groups []domain.StudentGroup) error {
	return g.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).Create(&groups).Error
}

func (g GroupDb) GetByID(ctx context.Context, groupID uint) (domain.StudentGroup, error) {
	var group domain.StudentGroup
	if err := g.db.WithContext(ctx).First(&group, groupID).Error; err != nil {
		return domain.StudentGroup{}, err
	}
	return group, nil
}

func (g GroupDb) GetByCampusID(ctx context.Context, campusID uint) ([]domain.StudentGroup, error) {
	var groups []domain.StudentGroup
	err := g.db.WithContext(ctx).Where("campus_id = ?", campusID).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (g GroupDb) GetByCampusIDs(ctx context.Context, campusIDs []uint) ([]domain.StudentGroup, error) {
	var groups []domain.StudentGroup
	err := g.db.WithContext(ctx).Where("campus_id IN ?", campusIDs).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}
func (g GroupDb) GetByCampusIDAndName(ctx context.Context, campusID uint, name string) ([]domain.StudentGroup, error) {
	var groups []domain.StudentGroup
	err := g.db.WithContext(ctx).Where("campus_id = ?", campusID).Where("name = ?", name).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}
func (g GroupDb) GetByCollegeIDAndName(ctx context.Context, collegeID uint, name string) ([]domain.StudentGroup, error) {
	var groups []domain.StudentGroup
	err := g.db.WithContext(ctx).Where("college_id = ?", collegeID).Where("name = ?", name).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return groups, nil
}
func (g GroupDb) GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.StudentGroup, error) {
	var groups []domain.StudentGroup
	if err := g.db.WithContext(ctx).
		Joins("JOIN campuses ON campuses.id = student_groups.campus_id").
		Where("campuses.college_id = ?", collegeID).
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
