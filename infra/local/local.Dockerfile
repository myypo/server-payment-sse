FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./internal/main.go

FROM alpine:3.20

WORKDIR /

COPY --from=build app/main /main
COPY --from=build app/.env /.env
COPY --from=build app/migration /migration

CMD ["./main"]
