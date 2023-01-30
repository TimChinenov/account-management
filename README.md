### Notes
curl http://localhost:8080/users \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "timmy","password": "test-password" }'

curl http://localhost:8080/login \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"username": "timmy","password": "test-password" }'