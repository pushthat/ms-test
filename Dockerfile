FROM golang:1.12 as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./main.go

FROM alpine:3.9

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN update-ca-certificates

COPY --from=build /app/main /app/main
RUN chmod 777 /app/main

ENTRYPOINT ["/app/main"]
