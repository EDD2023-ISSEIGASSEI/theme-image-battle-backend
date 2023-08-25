FROM golang:1.19.2-alpine

WORKDIR /go/app
COPY . .

RUN apk upgrade --update && \
  apk --no-cache add git

RUN go install github.com/cosmtrek/air@latest

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD ["air", "-c", ".air.toml"]