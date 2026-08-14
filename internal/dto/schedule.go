package dto

import (
	"encoding/json"
	"strconv"
	"time"
)

type Lesson struct {
	LessonID       uint   `json:"lessonId"`
	Title          string `json:"title"`
	Cabinet        string `json:"cabinet"`
	Date           Date   `json:"date"`
	Teacher        string `json:"teacher"`
	Order          uint   `json:"order"`
	StudentGroupID uint   `json:"studentGroupID"`
}

type ScheduleLessonResponse struct {
	Title     string     `json:"title"`
	Cabinet   string     `json:"cabinet"`
	Teacher   string     `json:"teacher"`
	Order     uint       `json:"order"`
	StartTime HourMinute `json:"startTime"`
	EndTime   HourMinute `json:"endTime"`
}

type ScheduleResponse struct {
	GroupID uint                     `json:"groupId"`
	Date    time.Time                `json:"date"`
	Lessons []ScheduleLessonResponse `json:"lessons"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil

}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(time.DateOnly))
}
