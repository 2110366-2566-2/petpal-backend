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

# References

- [Gin Documentation](https://gin-gonic.com/docs/quickstart/)

- [Golang Mongodb driver](https://www.mongodb.com/docs/drivers/go/current/usage-examples/)
