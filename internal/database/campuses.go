package database

import "github.com/ThisIsHyum/OpenScheduleApi/internal/models"

func (db *Db) GetCampusByID(id uint) (models.Campus, error) {
	var campus models.Campus
	return campus, db.db.Preload("StudentGroups").Find(&campus, id).Error
}
func (db *Db) GetCampusesByName(collegeId uint, name string) ([]models.Campus, error) {
	var campuses []models.Campus
	return campuses, db.db.Preload("StudentGroups").Where(
		models.Campus{CollegeID: collegeId, Name: name},
	).Find(&campuses).Error
}
func (db *Db) GetCampusesByCollegeID(collegeID uint) ([]models.Campus, error) {
	var campuses []models.Campus
	return campuses, db.db.Preload("StudentGroups").
		Where("college_id = ?", collegeID).Find(&campuses).Error
}
