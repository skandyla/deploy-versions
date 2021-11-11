package internal

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/skandyla/deploy-versions/config"
	"github.com/skandyla/deploy-versions/models"
)

type VersionStorage struct {
	db *sqlx.DB
}

var (
	tableName   = "versions"
	getAllQuery = `SELECT project, env, region, service, git_branch, git_commit_hash, build_id, user_name, created_at
                    	FROM ` + tableName + `
                    	ORDER BY created_at DESC
                    	LIMIT $1`
)

func NewVersionStorage(cfg *config.Config) (*VersionStorage, error) {
	connection, err := sqlx.Connect("postgres", cfg.PostgresDSN)
	if err != nil {
		return nil, fmt.Errorf("connection is not initialized %w", err)
	}

	return &VersionStorage{
		db: connection,
	}, nil
}

// Health checks availability of storage
func (s *VersionStorage) Health() error {
	return s.db.Ping()
}

// GetAll return all entities
func (s *VersionStorage) GetAllVersions() (m []models.VersionDBModel, err error) {
	limit := 10 //TBD - set as variable
	err = s.db.Select(&m, getAllQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("VersionsStorage.GetAllVersions: cannot select row : %s ", err)
	}
	return m, nil
}

// GetVersion
func (s *VersionStorage) GetVersionByID(id int) (m models.VersionDBModel, err error) {
	err = s.db.QueryRow("select * from versions where build_id = $1", id).
		Scan(&m.Project, &m.Env, &m.Region, &m.Service, &m.GitBranch, &m.GitCommitHash, &m.BuildID, &m.UserName, &m.CreatedAt)
	return m, err
}

// PostVersion - create new entity
func (s *VersionStorage) PostVersion(m models.VersionDBModel) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("PostVersion: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.Exec("INSERT into versions(project, env, region, service, git_branch, git_commit_hash, build_id, user_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		&m.Project, &m.Env, &m.Region, &m.Service, &m.GitBranch, &m.GitCommitHash, &m.BuildID, &m.UserName)

	log.Println(res)
	if err != nil {
		return fmt.Errorf("PostVersion: %w", err)
	}

	return tx.Commit()
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
