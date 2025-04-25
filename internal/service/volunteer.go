package service

import (
	"book-management/internal/controller"
	"context"
)

type VolunteerRepo interface {
	CreateVolunteer(ctx context.Context, volunteer controller.VolunteerInfo) (uint64, error)
	QueryVolunteers(ctx context.Context, pageSize, page int, id *uint64, total *int) ([]controller.VolunteerInfo, error)
	GetApplications(ctx context.Context, pageSize, page int, total *int) ([]controller.VolunteerApplicationInfo, error)
}

type VolunteerSvc struct {
	repo VolunteerRepo
}

func NewVolunteerSvc(repo VolunteerRepo) *VolunteerSvc {
	return &VolunteerSvc{repo: repo}
}

func (s *VolunteerSvc) CreateVolunteer(ctx context.Context, req controller.CreateVolunteerReq) (uint64, error) {
	volunteer := controller.VolunteerInfo{
		Name:                  req.Name,
		Phone:                 req.Phone,
		Age:                   req.Age,
		ServiceTimePreference: req.ServiceTimePreference,
		ExpertiseArea:         req.ExpertiseArea,
	}

	return s.repo.CreateVolunteer(ctx, volunteer)
}

func (s *VolunteerSvc) GetVolunteerInfos(ctx context.Context, req controller.GetVolunteerInfosReq, totalNum *int) ([]controller.VolunteerInfo, error) {
	return s.repo.QueryVolunteers(ctx, req.PageSize, req.Page, req.ID, totalNum)
}

func (s *VolunteerSvc) GetApplications(ctx context.Context, req controller.GetVolunteerApplicationsReq, total *int) ([]controller.VolunteerApplicationInfo, error) {
	return s.repo.GetApplications(ctx, req.PageSize, req.Page, total)
}
