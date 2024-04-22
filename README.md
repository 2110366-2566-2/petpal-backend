# petpal-backend

# Installation

1. must install Go version >1.13
1. run following commands in terminal
   1. `cd petpal-backend` (if not already in petpal-backend directory)
   1. `go mod tidy` (to install all the dependacies needed)
1. install docker

# Run

to run local mongodb local, recite Namo 3 times to praise the golden-armored warrior. and type the command

```bash
make run-develop
## or if you cannot run make
docker-compose -f docker-compose.dev.yml up -d
```

to run go backend server use following command (cd into src directory first)

```bash
cd src
go run main.go
```
and what if u want to use hot-reload(nodemon) 
```bash
npm install
cd src
nodemon
```

# Swagger 

## Installation
```bash
# to install swagger
> go install github.com/go-swagger/go-swagger/cmd/swagger@latest
# or use this if you clone this repo.
> go mod tidy
```

## Re-generate swagger
```bash
> swag init --parseInternal --parseDependency  -g main.go
```

## Accessing Swagger
for local, you can access swagger at `http://localhost:8080/swagger/index.html`

# References

- [Gin Documentation](https://gin-gonic.com/docs/quickstart/)

- [Golang Mongodb driver](https://www.mongodb.com/docs/drivers/go/current/usage-examples/)

- [Swagger declarative comment](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format)