package internal

import (
	"fmt"

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
	getAllQuery = `SELECT project, env, region, service, git_branch, git_commit_hash, build_id, created_at, created_at
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
