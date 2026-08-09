package database

import (
	"context"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Call struct {
	ID        uint
	Weekday   time.Weekday
	Begins    datatypes.Time
	Ends      datatypes.Time
	Order     uint `gorm:"column:call_order"`
	CollegeID uint
}

func newCall(call domain.Call) Call {
	return Call{
		ID:        call.ID,
		Weekday:   call.Weekday,
		Begins:    toDatatypesTime(call.Begins),
		Ends:      toDatatypesTime(call.Ends),
		Order:     call.Order,
		CollegeID: call.CollegeID,
	}
}

func (m Call) ToDomain() domain.Call {
	return domain.Call{
		ID:        m.ID,
		Weekday:   m.Weekday,
		Begins:    timeFromDatatypes(m.Begins),
		Ends:      timeFromDatatypes(m.Ends),
		Order:     m.Order,
		CollegeID: m.CollegeID,
	}
}

func toDatatypesTime(t time.Time) datatypes.Time {
	return datatypes.Time(
		time.Duration(t.Hour())*time.Hour +
			time.Duration(t.Minute())*time.Minute,
	)
}
func timeFromDatatypes(t datatypes.Time) time.Time {
	dur := time.Duration(t)
	hours := int(dur / time.Hour)
	minutes := int(dur/time.Minute) % 60
	return time.Date(1970, time.January, 1, hours, minutes, 0, 0, time.UTC)
}

type CallDb struct {
	db *gorm.DB
}

func NewCallDb(db *gorm.DB) *CallDb {
	return &CallDb{db: db}
}

func (c CallDb) GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.Call, error) {
	var callModels []Call
	if err := c.db.WithContext(ctx).Where("college_id = ?", collegeID).Find(&callModels).Error; err != nil {
		return nil, err
	}
	calls := make([]domain.Call, 0, len(callModels))
	for _, call := range callModels {
		calls = append(calls, call.ToDomain())
	}
	return calls, nil
}

func (c CallDb) UpsertMany(ctx context.Context, calls []domain.Call) error {
	callModels := make([]Call, 0, len(calls))
	for _, call := range calls {
		callModels = append(callModels, newCall(call))
	}
	return c.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "college_id"}, {Name: "weekday"}, {Name: "call_order"}},
			DoUpdates: clause.AssignmentColumns([]string{"begins", "ends"}),
		},
	).Create(&callModels).Error
}
