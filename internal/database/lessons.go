package database

import (
	"context"
	"slices"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/database/models"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LessonDb struct {
	db *gorm.DB
}

func NewLessonDb(db *gorm.DB) *LessonDb {
	return &LessonDb{db: db}
}

func (l LessonDb) Add(ctx context.Context, lessons []domain.Lesson) error {
	lessons = slices.DeleteFunc(lessons, func(lesson domain.Lesson) bool {
		return lesson.StudentGroupID == 0 || lesson.Title == ""
	})
	modelsLessons := []models.Lesson{}
	for _, lesson := range lessons {
		modelsLessons = append(modelsLessons, models.Lesson{
			Title:          lesson.Title,
			Cabinet:        lesson.Cabinet,
			Date:           datatypes.Date(lesson.Date),
			Teacher:        lesson.Teacher,
			Order:          lesson.Order,
			StudentGroupID: lesson.StudentGroupID,
		})
	}
	return l.db.WithContext(ctx).
		Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&modelsLessons).Error
}

func (l LessonDb) GetForDate(ctx context.Context, groupID uint, date time.Time) ([]domain.Lesson, error) {
	var lessons []models.Lesson
	if err := l.db.WithContext(ctx).
		Where("date = ? AND student_group_id = ?",
			date.Format("2006-01-02"), groupID).Find(&lessons).Error; err != nil {
		return nil, err
	}
	var domainLessons []domain.Lesson
	for _, lesson := range lessons {
		domainLessons = append(domainLessons, lesson.ToDomain())
	}
	return domainLessons, nil
}

func (l LessonDb) GetForDates(ctx context.Context, groupID uint,
	start, end time.Time) ([]domain.Lesson, error) {
	var lessons []models.Lesson
	if err := l.db.WithContext(ctx).
		Where("date BETWEEN ? AND ? AND student_group_id = ?", start, end, groupID).
		Find(&lessons).Error; err != nil {
		return nil, err
	}
	var domainLessons []domain.Lesson
	for _, lesson := range lessons {
		domainLessons = append(domainLessons, lesson.ToDomain())
	}
	return domainLessons, nil
}
