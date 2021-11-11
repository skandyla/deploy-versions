# deploy-versions

Deploy-version Microservice helps Devops teams to maintain history of deployed application versions

### simple cli testing
`docker compose up`  
`export POSTGRES_DSN="postgres://db_user:db_pass@localhost/deploy_versions?sslmode=disable"`  
`go run main.go`  
`curl localhost:8080/info`  
`curl -v -X GET localhost:8080/versions|jq`  
