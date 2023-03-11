### Notes
This is not so much of a readme as a notes page and quick commands page.

curl http://localhost:8080/api/users \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "test34","password": "123" }'

curl http://localhost:8080/api/login \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "test32","password": "123" }'


curl http://localhost:8080/api/admin/user \
    --include \
	--header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Nzg1MjkwMzYsInVzZXJJZCI6NDB9.B5MPt8nmYnKwTsWIzFgNIBnN565uMberEz9WkyAkG1E" \
    --request "GET"

curl http://localhost:8080/api/admin/posts/upvote \
    --include \
	--header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2Nzg1MjkwMzYsInVzZXJJZCI6NDB9.B5MPt8nmYnKwTsWIzFgNIBnN565uMberEz9WkyAkG1E" \
    --request "POST" \
    --data '{"userId": 40, "postId": 1, "voteType": 0}'

curl http://localhost:8080/api/admin/posts --include --header "Authorization: Bearer <token>" --header "Content-Type: application/json" --request "POST" --data '{ "userId": 26, "body": "lorem ipsum this is a test comment" }'

curl -X GET "http://localhost:8080/api/admin/posts/1/10" --include --header "Authorization: Bearer <token>"

### Useful Sources
https://github.com/gothinkster/golang-gin-realworld-example-app

https://github.com/Qovery/simple-example-gin-with-postgresql/blob/master/db/db.go

Notes on dependency inject:
https://stackoverflow.com/questions/46141898/what-is-the-best-way-to-have-dependency-injection-in-golang

### Guide on session handling
JWT explanation
https://seefnasrul.medium.com/create-your-first-go-rest-api-with-jwt-authentication-in-gin-framework-dbe5bda72817


### PSQL connection
psql -h "localhost" -U "postgres" -p "15432" -d "postgres" -a -f "/Users/timchinenov/Dev/account-management/migrations/2_create_posts_table.sql"

steps to migrate fly.io
    - proxy into server
    - run migrations with command above
    can we streamline this?