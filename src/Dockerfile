FROM golang:1.17-alpine3.14 as builder

WORKDIR $GOPATH/src/houze_ops_backend
COPY . .

RUN apk add --no-cache make curl git build-base

ENV GO111MODULE=on
ENV GIN_MODE=release

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /api main.go

FROM golang:1.17-alpine3.14

EXPOSE 8080

COPY --from=builder /api .

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN ["chmod", "+x", "/docker-entrypoint.sh"]