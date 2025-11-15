package database

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"gorm.io/gorm/clause"
)

func (db *Db) UpdateGroups(campusID uint, groupNames []string) error {
	var groups = []models.StudentGroup{}
	for _, groupName := range groupNames {
		groups = append(groups, models.StudentGroup{
			Name:     groupName,
			CampusID: campusID,
		})
	}
	return db.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&groups).Error
}

func (db *Db) GetGroupsByCampusID(campusID uint) ([]models.StudentGroup, error) {
	var groups []models.StudentGroup
	return groups, db.db.Where("campus_id = ?", campusID).Find(&groups).Error
}
func (db *Db) GetGroupByID(ID uint) (models.StudentGroup, error) {
	var group models.StudentGroup
	return group, db.db.First(&group, ID).Error
}
func (db *Db) GetGroupsByName(campusID uint, name string) ([]models.StudentGroup, error) {
	var groups []models.StudentGroup
	return groups, db.db.Where(&models.StudentGroup{
		CampusID: campusID, Name: name}).
		Find(&groups).Error
}

func (db *Db) GetGroupsByCollegeID(collegeID uint) ([]models.StudentGroup, error) {
	var groups []models.StudentGroup
	return groups, db.db.
		Joins("JOIN campuses ON campuses.campus_id = student_groups.campus_id").
		Where("campuses.college_id = ?", collegeID).
		Find(&groups).Error
}
