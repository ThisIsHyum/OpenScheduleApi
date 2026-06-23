package models

import (
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"gorm.io/datatypes"
)

type (
	College struct {
		CollegeID uint     `gorm:"primaryKey;unique;autoIncrement"`
		Name      string   `gorm:"unique" json:"name"`
		Calls     []Call   `gorm:"constraint:OnDelete:CASCADE;"`
		Campuses  []Campus `gorm:"constraint:OnDelete:CASCADE;"`
		Token     string   `gorm:"token"`
	}

	Campus struct {
		CampusID      uint           `gorm:"primaryKey;unique;autoIncrement" json:"campusId"`
		Name          string         `json:"name"`
		CollegeID     uint           `json:"collegeId"`
		StudentGroups []StudentGroup `gorm:"constraint:OnDelete:CASCADE;" json:"studentGroups,omitempty"`
	}

	StudentGroup struct {
		StudentGroupID uint   `gorm:"primaryKey;autoIncrement" json:"studentGroupId,omitempty"`
		Name           string `gorm:"type:varchar(100);not null;uniqueIndex:idx_campus_name"`
		CampusID       uint   `gorm:"not null;uniqueIndex:idx_campus_name"`
	}

	Lesson struct {
		LessonID       uint           `gorm:"primaryKey;autoIncrement" json:"lessonId"`
		Title          string         `json:"title"`
		Cabinet        string         `json:"cabinet"`
		Date           datatypes.Date `gorm:"uniqueIndex:idx_group_date_order" json:"date"`
		Teacher        string         `json:"teacher"`
		Order          uint           `gorm:"uniqueIndex:idx_group_date_order" json:"order"`
		StudentGroupID uint           `gorm:"uniqueIndex:idx_group_date_order" json:"studentGroupID"`
		StudentGroup   StudentGroup   `gorm:"constraint:OnDelete:CASCADE;" json:"studentGroup"`
	}

	Call struct {
		CallID    uint           `gorm:"primaryKey;unique;autoIncrement" json:"callId"`
		Weekday   time.Weekday   `gorm:"uniqueIndex:idx_weekday_college_order" json:"weekday"`
		Begins    datatypes.Time `json:"begins"`
		Ends      datatypes.Time `json:"ends"`
		Order     uint           `gorm:"uniqueIndex:idx_weekday_college_order" json:"order"`
		CollegeID uint           `gorm:"uniqueIndex:idx_weekday_college_order" json:""`
	}
)

func (Campus) TableName() string { return "campuses" }

func (l Lesson) ToDomain() domain.Lesson {
	return domain.Lesson{
		ID:             l.LessonID,
		Title:          l.Title,
		Cabinet:        l.Cabinet,
		Date:           time.Time(l.Date),
		Teacher:        l.Teacher,
		Order:          l.Order,
		StudentGroupID: l.StudentGroupID,
	}
}
func (c College) ToDomain() domain.College {
	return domain.College{
		ID:   c.CollegeID,
		Name: c.Name,
	}
}

func (c Campus) ToDomain() domain.Campus {
	return domain.Campus{
		ID:        c.CampusID,
		Name:      c.Name,
		CollegeID: c.CollegeID,
	}
}
func (sg StudentGroup) ToDomain() domain.StudentGroup {
	return domain.StudentGroup{
		ID:       sg.StudentGroupID,
		Name:     sg.Name,
		CampusID: sg.CampusID,
	}
}

func (c Call) ToDomain() domain.Call {
	return domain.Call{
		ID:        c.CallID,
		Weekday:   c.Weekday,
		Begins:    timeFromDatatypes(c.Begins),
		Ends:      timeFromDatatypes(c.Ends),
		Order:     c.Order,
		CollegeID: c.CollegeID,
	}
}

func timeFromDatatypes(t datatypes.Time) time.Time {
	dur := time.Duration(t)
	hours := int(dur / time.Hour)
	minutes := int(dur/time.Minute) % 60
	return time.Date(1970, time.January, 1, hours, minutes, 0, 0, time.UTC)
}
