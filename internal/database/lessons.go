package database

import (
	"slices"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/models"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
)

func (db *Db) AddLessons(lessons []models.Lesson) error {
	lessons = slices.DeleteFunc(lessons, func(lesson models.Lesson) bool {
		return lesson.StudentGroupID == 0 || lesson.Title == ""
	})
	return db.db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&lessons).Error
}

func (db *Db) GetLessonsForDate(groupID uint, date time.Time) ([]models.Lesson, error) {
	var lessons []models.Lesson
	return lessons, db.db.Where("date = ? AND student_group_id = ?",
		date.Format("2006-01-02"), groupID).Find(&lessons).Error
}

func (db *Db) GetLessonsForDates(groupID uint, start, end datatypes.Date) ([]models.Lesson, error) {
	var lessons []models.Lesson
	return lessons, db.db.
		Where("date BETWEEN ? AND ? AND student_group_id = ?", start, end, groupID).
		Find(&lessons).Error
}
