check available flags: `go run ./cmd/api --help`

Golang version: go1.22.2
backend start: port 4000

1. `cd backend`
2. `go mod tidy`
3. `go run ./cmd/api`

node version: v22.15.0
frontend start: port 3000

1. `cd frontend`
2. `npm i && npm start`

curl:

> get: curl -w '\nTime: %{time_total}s \n' localhost:4000/contacts
> post: curl -i -d "{"first_name": "myname", "last_name":"familyname", "email": "test@test.com", "phone_number":"1231231234"}" localhost:4000/contacts
> patch: curl -X PATCH -d '{"email":"new_test_email@test.com"}' localhost:4000/contacts/1
> delete: curl -X DELETE localhost:4000/contacts/1
