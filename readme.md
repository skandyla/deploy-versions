# deploy-versions

Deploy-version Microservice helps Devops teams to maintain history of deployed application versions

### simple cli testing
`docker compose up`  
`export POSTGRES_DSN="postgres://db_user:db_pass@localhost/deploy_versions?sslmode=disable"`  
`go run main.go`  
`curl localhost:8080/info`  
`curl -v -X GET localhost:8080/versions|jq`  
`curl -v localhost:8080/version/130`  
`curl -v -X POST localhost:8080/version -d '{"project":"example","env":"dev","region":"us-west-1","service_name"="notifications","user_name":"bsmith","build_id"=133}'`  
`curl -v localhost:8080/version/133`  
`curl -v localhost:8080/version/140`  #not found expected  