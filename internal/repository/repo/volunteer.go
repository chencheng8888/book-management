package repo

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"context"
)

type VolunteerDao interface {
	Create(ctx context.Context, volunteer do.Volunteer) error
	Query(ctx context.Context, pageSize, page int, id *uint64) ([]do.Volunteer, error)
	Count(ctx context.Context, id *uint64) (int, error)
	QueryApplications(ctx context.Context, pageSize, page int) ([]do.VolunteerApplication, error)
	CountApplications(ctx context.Context) (int, error)
}

type VolunteerRepo struct {
	dao VolunteerDao
}

func NewVolunteerRepo(dao VolunteerDao) *VolunteerRepo {
	return &VolunteerRepo{dao: dao}
}

func (r *VolunteerRepo) CreateVolunteer(ctx context.Context, volunteer controller.VolunteerInfo) (uint64, error) {
	volunteerDO := do.Volunteer{
		Name:                  volunteer.Name,
		Phone:                 volunteer.Phone,
		Age:                   volunteer.Age,
		ServiceTimePreference: volunteer.ServiceTimePreference,
		ExpertiseArea:         volunteer.ExpertiseArea,
	}

	if err := r.dao.Create(ctx, volunteerDO); err != nil {
		return 0, errcode.AddDataError
	}

	return volunteerDO.ID, nil
}

func (r *VolunteerRepo) QueryVolunteers(ctx context.Context, pageSize, page int, id *uint64, total *int) ([]controller.VolunteerInfo, error) {
	count, err := r.dao.Count(ctx, id)
	if err != nil {
		return nil, err
	}
	*total = count

	if maxPage := tool.GetPage(count, pageSize); page > maxPage {
		return nil, errcode.PageError
	}

	volunteers, err := r.dao.Query(ctx, pageSize, page, id)
	if err != nil {
		return nil, err
	}

	return batchToVolunteerController(volunteers), nil
}

func (r *VolunteerRepo) GetApplications(ctx context.Context, pageSize, page int, total *int) ([]controller.VolunteerApplicationInfo, error) {
	count, err := r.dao.CountApplications(ctx)
	if err != nil {
		return nil, err
	}
	*total = count

	applications, err := r.dao.QueryApplications(ctx, pageSize, page)
	if err != nil {
		return nil, err
	}

	res := make([]controller.VolunteerApplicationInfo, 0, len(applications))
	for _, app := range applications {
		res = append(res, controller.VolunteerApplicationInfo{
			ID:    app.ID,
			Name:  app.Name,
			Phone: app.Phone,
			Age:   app.Age,
		})
	}
	return res, nil
}

func batchToVolunteerController(dos []do.Volunteer) []controller.VolunteerInfo {
	res := make([]controller.VolunteerInfo, 0, len(dos))
	for _, d := range dos {
		res = append(res, controller.VolunteerInfo{
			ID:                    d.ID,
			Name:                  d.Name,
			Phone:                 d.Phone,
			Age:                   d.Age,
			ServiceTimePreference: d.ServiceTimePreference,
			ExpertiseArea:         d.ExpertiseArea,
			CreatedAt:             d.CreatedAt.Format("2006-01-02"),
		})
	}
	return res
}
