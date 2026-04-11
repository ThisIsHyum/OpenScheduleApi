package service

import (
	"context"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/domain"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/repository"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/token"
)

type AdminService struct {
	collegeRepo repository.CollegeRepo
	campusRepo  repository.CampusRepo
	createTx    repository.CreateTx
}

func NewAdminService(collegeRepo repository.CollegeRepo,
	campusRepo repository.CampusRepo,
	createTx repository.CreateTx) *AdminService {
	return &AdminService{
		collegeRepo: collegeRepo, campusRepo: campusRepo, createTx: createTx,
	}
}

func (s AdminService) NewParser(ctx context.Context,
	collegeName string, campusNames []string) (string, error) {
	if colleges, err := s.collegeRepo.GetByName(ctx, collegeName); err != nil {
		return "", err
	} else if len(colleges) != 0 {
		return "", domain.ErrConflict
	}
	tx, err := s.createTx()
	if err != nil {
		return "", err
	}

	defer tx.Rollback()

	collegeRepo := s.collegeRepo.WithTx(tx)
	campusRepo := s.campusRepo.WithTx(tx)

	token, err := token.GenerateToken()
	if err != nil {
		return "", err
	}

	collegeID, err := collegeRepo.Create(ctx, domain.College{
		Name:  collegeName,
		Token: token})
	if err != nil {
		return "", err
	}

	campuses := make([]domain.Campus, len(campusNames))
	for i, campusName := range campusNames {
		campuses[i] = domain.Campus{Name: campusName, CollegeID: collegeID}
	}
	if err := campusRepo.CreateMany(ctx, campuses); err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return token, nil
}

func (s AdminService) DeleteParser(ctx context.Context, parserID uint) error {
	return s.collegeRepo.Delete(ctx, parserID)
}
