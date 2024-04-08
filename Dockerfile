FROM golang:1.18-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
WORKDIR ./src

# Command to run the executable
CMD ["go","run","main.go"]
# CMD ["bin bash","echo","Hello World!"]