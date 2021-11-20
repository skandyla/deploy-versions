package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/skandyla/deploy-versions/models"
)

type VersionRepository struct {
	db *sqlx.DB
}

var (
	tableName   = "versions"
	getAllQuery = `SELECT project, env, region, service, git_branch, git_commit_hash, build_id, user_name, created_at
                    	FROM ` + tableName + `
                    	ORDER BY created_at DESC
                    	LIMIT $1`
)

func NewVersionRepository(db *sqlx.DB) *VersionRepository {
	return &VersionRepository{db: db}
}

// Health checks availability of storage
func (s *VersionRepository) Health(ctx context.Context) error {
	return s.db.Ping()
}

// GetAll return all entities
func (s *VersionRepository) GetAll(ctx context.Context) (m []models.VersionResponse, err error) {
	limit := 10 //TBD - set as variable
	dbModels := []models.VersionDBModel{}
	err = s.db.Select(&dbModels, getAllQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("VersionsRepository.GetAll: cannot select row : %s ", err)
	}

	for _, v := range dbModels {
		m = append(m, versionResponseFromDBModel(v))
	}

	return m, nil
}

// GetVersion
func (s *VersionRepository) GetByID(ctx context.Context, id int) (m models.VersionDBModel, err error) {
	err = s.db.QueryRow("select * from versions where build_id = $1", id).
		Scan(&m.Project, &m.Env, &m.Region, &m.Service, &m.GitBranch, &m.GitCommitHash, &m.BuildID, &m.UserName, &m.CreatedAt)
	return m, err
}

func (s *VersionRepository) Delete(ctx context.Context, id int) error {
	//TBD delete
	log.Println("Deleting ID:", id)
	return nil
}

// PostVersion - create new entity
func (s *VersionRepository) Create(ctx context.Context, version models.CreateVersionRequest) error {
	m := versionDBModelFromCreateRequest(version)
	//tx, err := s.db.Begin()
	//if err != nil {
	//	return fmt.Errorf("PostVersion: %w", err)
	//}
	//defer tx.Rollback() //need to check error

	res, err := s.db.Exec("INSERT into versions(project, env, region, service, git_branch, git_commit_hash, build_id, user_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		&m.Project, &m.Env, &m.Region, &m.Service, &m.GitBranch, &m.GitCommitHash, &m.BuildID, &m.UserName)

	log.Println(res)
	if err != nil {
		return fmt.Errorf("PostVersion: %w", err)
	}
	return nil
	//return tx.Commit()
}

//translate DBModel into our ResponseStruct
func versionResponseFromDBModel(v models.VersionDBModel) models.VersionResponse {
	return models.VersionResponse{
		Env:       v.Env,
		Region:    v.Region,
		Service:   v.Service,
		UserName:  v.UserName,
		CreatedAt: v.CreatedAt,
	}
}

//for creating new versions - we can accept different representation of object, than stored in our DB for single entity
func versionDBModelFromCreateRequest(r models.CreateVersionRequest) models.VersionDBModel {
	return models.VersionDBModel{
		Project:       r.Project,
		Env:           r.Env,
		Region:        r.Region,
		Service:       r.Service,
		GitBranch:     r.GitBranch,
		GitCommitHash: r.GitCommitHash,
		BuildID:       r.BuildID,
		UserName:      r.UserName,
	}
}
