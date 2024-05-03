FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/blockd cmd/main.go

EXPOSE 8080

CMD ["/app/blockd", "-log-level=info", "-address=0.0.0.0:8080", "-db-host=blockd-db:8432", "-db-user=blockd", "-db-password=blockd"]
