FROM golang:1.19.2-alpine

WORKDIR /go/app
COPY ./app .

RUN apk upgrade --update && \
  apk --no-cache add git

RUN go get -u github.com/cosmtrek/air && \
  go build -o /go/bin/air github.com/cosmtrek/air

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

CMD ["air", "-c", ".air.toml"]