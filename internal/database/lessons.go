package database

import (
	"context"
	"slices"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Lesson struct {
	ID             uint
	Title          string
	Cabinet        string
	Teacher        string
	Date           time.Time
	Order          uint `gorm:"column:lesson_order"`
	StudentGroupID uint
}

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

	lessonModels := make([]Lesson, 0, len(lessons))
	for _, lesson := range lessons {
		lessonModels = append(lessonModels, Lesson(lesson))
	}
	return l.db.WithContext(ctx).
		Clauses(clause.Insert{Modifier: "IGNORE"}).Create(&lessonModels).Error
}

func (l LessonDb) GetForDate(ctx context.Context, groupID uint, date time.Time) ([]domain.Lesson, error) {
	var lessonModels []Lesson
	if err := l.db.WithContext(ctx).
		Where("date = ? AND student_group_id = ?",
			date.Format("2006-01-02"), groupID).Find(&lessonModels).Error; err != nil {
		return nil, err
	}
	lessons := make([]domain.Lesson, 0, len(lessonModels))
	for _, lesson := range lessonModels {
		lessons = append(lessons, domain.Lesson(lesson))
	}
	return lessons, nil
}

func (l LessonDb) GetForDates(ctx context.Context, groupID uint,
	start, end time.Time) ([]domain.Lesson, error) {
	var lessonModels []Lesson
	if err := l.db.WithContext(ctx).
		Where("date BETWEEN ? AND ? AND student_group_id = ?",
			start.Format("2006-01-02"), end.Format("2006-01-02"), groupID).
		Find(&lessonModels).Error; err != nil {
		return nil, err
	}
	lessons := make([]domain.Lesson, 0, len(lessonModels))
	for _, lesson := range lessonModels {
		lessons = append(lessons, domain.Lesson(lesson))
	}
	return lessons, nil
}
