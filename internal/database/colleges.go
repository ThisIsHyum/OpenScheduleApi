package database

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
)

func (db Db) GetColleges() ([]models.College, error) {
	var colleges []models.College
	err := db.db.Preload("Calls").Preload("Campuses").Find(&colleges).Error
	return colleges, err
}

func (db Db) GetCollege(ID uint) (models.College, error) {
	var college models.College
	return college, db.db.Preload("Calls").Preload("Campuses").
		First(&college, ID).Error
}

func (db Db) GetCollegesByName(name string) ([]models.College, error) {
	colleges := []models.College{}
	return colleges, db.db.Where("name = ?", name).Preload("Calls").
		Preload("Campuses").Find(&colleges).Error
}

func (db Db) GetCollegeIDByGroupID(groupID uint) (uint, error) {
	var collegeID uint
	return collegeID, db.db.Table("student_groups AS sg").
		Select("c.college_id").
		Joins("JOIN campuses AS c ON c.campus_id = sg.campus_id").
		Where("sg.student_group_id = ?", groupID).
		Scan(&collegeID).Error
}
