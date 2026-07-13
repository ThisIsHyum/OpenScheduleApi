package repository

import (
	"context"
	"time"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
)

type CollegeRepo interface {
	Create(ctx context.Context, college domain.College) (uint, error)
	Delete(ctx context.Context, parserId uint) error
	Get(ctx context.Context, collegeID uint) (domain.College, error)
	GetByName(ctx context.Context, name string) ([]domain.College, error)
	GetByToken(ctx context.Context, token string) (domain.College, error)
	GetAll(ctx context.Context) ([]domain.College, error)
	GetByGroupID(ctx context.Context, groupID uint) (domain.College, error)

	WithTx(tx Tx) CollegeRepo
}

type CallRepo interface {
	UpsertMany(ctx context.Context, calls []domain.Call) error
	GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.Call, error)
}

type StudentGroupRepo interface {
	UpsertMany(ctx context.Context, groups []domain.StudentGroup) error
	GetByID(ctx context.Context, groupID uint) (domain.StudentGroup, error)
	GetByCampusID(ctx context.Context, campusID uint) ([]domain.StudentGroup, error)
	GetByCampusIDs(ctx context.Context, campusIDs []uint) ([]domain.StudentGroup, error)
	GetByCampusIDAndName(ctx context.Context, campusID uint, name string) ([]domain.StudentGroup, error)
	GetByCollegeIDAndName(ctx context.Context, collegeID uint, name string) ([]domain.StudentGroup, error)
	GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.StudentGroup, error)
}

type LessonRepo interface {
	Add(ctx context.Context, lessons []domain.Lesson) error
	GetForDate(ctx context.Context, groupID uint, date time.Time) ([]domain.Lesson, error)
	GetForDates(ctx context.Context, groupID uint, start, end time.Time) ([]domain.Lesson, error)
}

type CampusRepo interface {
	CreateMany(ctx context.Context, campuses []domain.Campus) error
	GetByID(ctx context.Context, campusID uint) (domain.Campus, error)
	GetByName(ctx context.Context, collegeId uint, name string) ([]domain.Campus, error)
	GetByCollegeID(ctx context.Context, collegeID uint) ([]domain.Campus, error)
	GetByCollegeIDs(ctx context.Context, collegeIDs []uint) ([]domain.Campus, error)

	WithTx(tx Tx) CampusRepo
}

type Tx interface {
	Commit() error
	Rollback() error
	DB() any
}

type CreateTx func() (Tx, error)

type Repos struct {
	CollegeRepo      CollegeRepo
	CallRepo         CallRepo
	StudentGroupRepo StudentGroupRepo
	LessonRepo       LessonRepo
	CampusRepo       CampusRepo
	CreateTx         CreateTx
}
