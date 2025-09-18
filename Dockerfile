FROM golang:1.24

WORKDIR /app

RUN go install github.com/air-verse/air@v1.62.0 && go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

COPY go.mod go.sum ./

RUN go mod download

COPY . .