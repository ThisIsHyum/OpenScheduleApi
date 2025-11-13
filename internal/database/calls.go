package database

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"gorm.io/gorm/clause"
)

func (db *Db) UpdateCalls(collegeID uint, calls []models.Call) error {
	for i := range calls {
		calls[i].CollegeID = collegeID
	}
	return db.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "college_id"}, {Name: "weekday"}, {Name: "order"}},
			DoUpdates: clause.AssignmentColumns([]string{"begins", "ends"}),
		},
	).Create(&calls).Error
}

func (db *Db) GetCalls(collegeId uint) ([]models.Call, error) {
	var calls []models.Call
	return calls, db.db.Where("college_id = ?", collegeId).Find(&calls).Error
}
