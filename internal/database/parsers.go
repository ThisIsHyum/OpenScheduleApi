package database

import (
	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
)

func (db *Db) NewParser(token string, collegeName string, campusNames []string) error {
	tx := db.db.Begin()

	college := models.College{Name: collegeName}
	parser := models.Parser{Token: token, College: college}
	if err := tx.Create(&parser).Error; err != nil {
		tx.Rollback()
		return err
	}

	campuses := []models.Campus{}
	for _, campusName := range campusNames {
		campuses = append(campuses, models.Campus{
			Name:      campusName,
			CollegeID: parser.College.CollegeID,
		})
	}
	if err := tx.Create(&campuses).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (db *Db) DeleteParser(parserId uint) (bool, error) {
	result := db.db.Delete(&models.Parser{}, parserId)
	return result.RowsAffected != 0, result.Error
}

func (db *Db) GetParserByToken(token string) (models.Parser, error) {
	var parser models.Parser
	return parser, db.db.Preload("College").Where("token = ?", token).First(&parser).Error
}
