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
