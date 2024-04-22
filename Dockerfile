# Start from the official Golang base image
FROM golang:latest AS builder
ENV GIN_MODE=release
WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/


FROM alpine:latest  
ENV GIN_MODE=release
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY ./src/.env ./src/.env
CMD ["./main"]

# FROM golang:latest AS builder

# WORKDIR /app

# COPY . ./
# RUN go mod download

# RUN CGO_ENABLED=0 go build -o /bin/app

# FROM debian:buster-slim

# COPY --from=build /bin/app /bin
# COPY ./src/.env ./src/.env

# EXPOSE 8080

# CMD [ "/bin/app" ]