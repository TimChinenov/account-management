### Notes
This is not so much of a readme as a notes page and quick commands page.

curl http://localhost:8080/api/users \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "test6","password": "123" }'

curl http://localhost:8080/api/login \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "test6","password": "123" }'


curl http://localhost:8080/api/admin/user \
    --include \
	--header "Authorization: Bearer <token>"
    --request "GET"

### Useful Sources
https://github.com/gothinkster/golang-gin-realworld-example-app

https://github.com/Qovery/simple-example-gin-with-postgresql/blob/master/db/db.go

Notes on dependency inject:
https://stackoverflow.com/questions/46141898/what-is-the-best-way-to-have-dependency-injection-in-golang

### Guide on session handling
JWT explanation
https://seefnasrul.medium.com/create-your-first-go-rest-api-with-jwt-authentication-in-gin-framework-dbe5bda72817


### PSQL connection
psql -h "localhost" -U "postgres" -p "15432" -d "postgres" -a -f "/Users/timchinenov/Dev/account-management/migrations/1_create_user_table.sql"

steps to migrate fly.io
    - proxy into server
    - run migrations with command above
    can we streamline this?