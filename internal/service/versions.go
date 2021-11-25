package service

import (
	"context"

	"github.com/skandyla/deploy-versions/models"
)

type VersionsService interface {
	Create(ctx context.Context, version models.CreateVersionRequest) error
	GetByID(ctx context.Context, id int) (models.VersionDBModel, error)
	GetAll(ctx context.Context) ([]models.VersionResponse, error)
	Delete(ctx context.Context, id int) error
	Health(ctx context.Context) error
	//Update(ctx context.Context, id int64, req models.UpdateVersionRequest) error
}

type Versions struct {
	repo VersionsService
}

func NewVersions(repo VersionsService) *Versions {
	return &Versions{
		repo: repo,
	}
}

func (v *Versions) Create(ctx context.Context, version models.CreateVersionRequest) error {
	//any custom app logic here
	if version.Region == "" {
		version.Region = "Global"
	}

	return v.repo.Create(ctx, version)
}

func (v *Versions) GetByID(ctx context.Context, id int) (models.VersionDBModel, error) {
	return v.repo.GetByID(ctx, id)
}

func (v *Versions) GetAll(ctx context.Context) ([]models.VersionResponse, error) {
	return v.repo.GetAll(ctx)
}

func (v *Versions) Delete(ctx context.Context, id int) error {
	return v.repo.Delete(ctx, id)
}

func (v *Versions) Health(ctx context.Context) error {
	return v.repo.Health(ctx)
}

//func (v *Versions) Update(ctx context.Context, id int64, inp models.UpdateVersionRequest) error {
//	return v.repo.Update(ctx, id, inp)
//}
