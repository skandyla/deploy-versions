package models

import "time"

type VersionDBModel struct {
	Project       string    `db:"project"`
	Env           string    `db:"env"`
	Region        string    `db:"region"`
	Service       string    `db:"service"`
	GitBranch     string    `db:"git_branch"`
	GitCommitHash string    `db:"git_commit_hash"`
	BuildID       string    `db:"build_id"`
	UserName      string    `db:"user_name"`
	CreatedAt     time.Time `db:"created_at"`
}

type VersionResponse struct {
	Env       string    `json:"env"`
	Region    string    `json:"region"`
	Service   string    `json:"service_name"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateVersionRequest struct {
	Project       string `json:"project"`
	Env           string `json:"env" validate:"required"`
	Region        string `json:"region" validate:"required"`
	Service       string `json:"service_name" validate:"required"`
	UserName      string `json:"user_name" validate:"required"`
	GitBranch     string `json:"git_branch"`
	GitCommitHash string `json:"git_commit_hash"`
	BuildID       string `json:"build_id" validate:"required"`
}
