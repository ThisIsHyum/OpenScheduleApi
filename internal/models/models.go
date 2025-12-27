package models

import (
	"time"

	"gorm.io/datatypes"
)

type (
	College struct {
		CollegeID uint     `gorm:"primaryKey;unique;autoIncrement" json:"collegeId"`
		Name      string   `gorm:"unique" json:"name"`
		Calls     []Call   `gorm:"constraint:OnDelete:CASCADE;" json:"calls"`
		Campuses  []Campus `gorm:"constraint:OnDelete:CASCADE;" json:"campuses"`
		ParserID  uint     `json:"-"`
	}

	Campus struct {
		CampusID      uint           `gorm:"primaryKey;unique;autoIncrement" json:"campusId"`
		Name          string         `json:"name"`
		CollegeID     uint           `json:"collegeId"`
		StudentGroups []StudentGroup `gorm:"constraint:OnDelete:CASCADE;" json:"studentGroups,omitempty"`
	}

	StudentGroup struct {
		StudentGroupID uint   `gorm:"primaryKey;autoIncrement" json:"studentGroupId,omitempty"`
		Name           string `gorm:"type:varchar(100);not null;uniqueIndex:idx_campus_name" json:"name"`
		CampusID       uint   `gorm:"not null;uniqueIndex:idx_campus_name" json:"campusId"`
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
	Parser struct {
		ParserID uint `gorm:"primaryKey;unique;autoIncrement"`
		Token    string
		College  College `gorm:"constraint:OnDelete:CASCADE"`
	}

	Schedule struct {
		GroupID uint             `json:"groupId"`
		Date    datatypes.Date   `json:"date"`
		Lessons []ScheduleLesson `json:"lessons"`
	}

	ScheduleLesson struct {
		Title     string         `json:"title"`
		Cabinet   string         `json:"cabinet"`
		Teacher   string         `json:"teacher"`
		Order     uint           `json:"order"`
		StartTime datatypes.Time `json:"startTime"`
		EndTime   datatypes.Time `json:"endTime"`
	}
)

func (Campus) TableName() string { return "campuses" }

func (c College) String() string      { return c.Name }
func (c Campus) String() string       { return c.Name }
func (g StudentGroup) String() string { return g.Name }

func findByDate(schedules []Schedule, date datatypes.Date) (Schedule, int, bool) {
	for i, schedule := range schedules {
		if schedule.Date == date {
			return schedule, i, true
		}
	}
	return Schedule{}, -1, false
}
func NewSchedule(lessons []Lesson, calls []Call) Schedule {
	schedule := Schedule{
		GroupID: lessons[0].StudentGroupID,
		Date:    lessons[0].Date,
		Lessons: []ScheduleLesson{},
	}
	for _, lesson := range lessons {
		value, err := lesson.Date.Value()
		if err != nil {
			continue
		}
		schedule.addLesson(lesson, calls, value.(time.Time).Weekday())
	}
	return schedule
}

func NewSchedules(lessons []Lesson, calls []Call) []Schedule {
	schedules := []Schedule{}
	for _, lesson := range lessons {
		_, i, exists := findByDate(schedules, lesson.Date)
		if exists {
			value, err := lesson.Date.Value()
			if err != nil {
				continue
			}
			schedules[i].addLesson(lesson, calls, value.(time.Time).Weekday())
		} else {
			schedule := Schedule{
				GroupID: lesson.StudentGroupID,
				Date:    lesson.Date,
				Lessons: []ScheduleLesson{},
			}
			value, err := lesson.Date.Value()
			if err != nil {
				continue
			}
			schedules[i].addLesson(lesson, calls, value.(time.Time).Weekday())
			schedules = append(schedules, schedule)
		}
	}
	return schedules
}

func (s *Schedule) addLesson(lesson Lesson, calls []Call, weekday time.Weekday) {
	s.Lessons = append(s.Lessons, lesson.toScheduleLesson(calls, weekday))
}

func (l Lesson) toScheduleLesson(calls []Call, weekday time.Weekday) ScheduleLesson {
	call := CallByOrderAndWeekday(calls, l.Order, weekday)
	return ScheduleLesson{
		Title:     l.Title,
		Cabinet:   l.Cabinet,
		Teacher:   l.Teacher,
		Order:     l.Order,
		StartTime: call.Begins,
		EndTime:   call.Ends,
	}
}

func CallByOrderAndWeekday(calls []Call, order uint, weekday time.Weekday) Call {
	for _, call := range calls {
		if call.Order == order && call.Weekday == weekday {
			return call
		}
	}
	return Call{}
}
