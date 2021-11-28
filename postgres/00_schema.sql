--CREATE DATABASE deploy_versions;
--\c deploy_versions

CREATE TABLE IF NOT EXISTS versions (
  project varchar(32),          -- Project name. We can have multiple
  env  varchar(32),             -- Environment name
  region  varchar(32),          -- Region were we deploy
  service  varchar(32),         -- Service name within a project
  git_branch  varchar(32),      -- Git branch/tag deployed from
  git_commit_hash  varchar(64), -- Git commit hash
  build_id  integer,            -- CI build_id
  user_name  varchar(32),       -- User that started deploy
  created_at timestamp NOT NULL DEFAULT now() -- Time of deploy
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL NOT NULL unique,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    registered_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL NOT NULL unique,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    token VARCHAR(255) NOT NULL unique, 
    expires_at TIMESTAMP NOT NULL
);


-- create record
INSERT INTO versions(project, env, region, service, git_branch, git_commit_hash, build_id, user_name, created_at) 
    VALUES('MyUnicorn1', 'stg', 'eu-central-1', 'api-gateway', 'v0.1.1', '7d0eb417009f5794a09330f0aad3934bea476a53', 129, 'jsmith' ,'2021-11-08 17:05:23.048055');
INSERT INTO versions(project, env, region, service, git_branch, git_commit_hash, build_id, user_name) 
    VALUES('MyUnicorn1', 'stg', 'eu-central-1', 'api-gateway', 'v0.1.1', '7d0eb417009f5794a09330f0aad3934bea476a53', 130, 'jsmith');
INSERT INTO versions(project, env, region, service, git_branch, git_commit_hash, build_id, user_name) 
    VALUES('MyUnicorn1', 'stg', 'eu-west-2', 'api-gateway', 'v0.1.1', '7d0eb417009f5794a09330f0aad3934bea476a53', 131, 'rpike');