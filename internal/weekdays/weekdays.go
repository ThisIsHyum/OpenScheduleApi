package weekdays

import (
	"strings"
	"time"

	"gorm.io/datatypes"
)

var stringsWeekdays = map[string]time.Weekday{}

func ParseWeekday(s string) (time.Weekday, bool) {
	w, ok := stringsWeekdays[strings.ToLower(s)]
	return w, ok
}

func init() {
	for i := range 7 {
		day := time.Weekday(i).String()
		stringsWeekdays[strings.ToLower(day)] = time.Weekday(i)
	}
}

func GetDateByWeekday(weekday time.Weekday, offset int) time.Time {
	now := time.Now().AddDate(0, 0, offset*7)
	current := now.Weekday()
	if current == 0 {
		current = 7
	}
	if weekday == 0 {
		weekday = 7
	}
	return now.AddDate(0, 0, int(weekday-current))
}

func WeekBounds(offset int) (datatypes.Date, datatypes.Date) {
	refDate := time.Now().AddDate(0, 0, offset*7)

	weekday := int(refDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	daysFromMonday := weekday - 1
	start := time.Date(
		refDate.Year(), refDate.Month(), refDate.Day(), 0, 0, 0, 0, refDate.Location(),
	).AddDate(0, 0, -daysFromMonday)
	end := start.AddDate(0, 0, 6)
	return datatypes.Date(start), datatypes.Date(end)
}
