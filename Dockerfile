# Heroku doesn't seem to fully support multi-stage builds yet
FROM golang:1.13.5-alpine

# Enable Go modules
ENV GO111MODULE=on
# Install Git so `go get` works
RUN apk add --no-cache git
# Install the Certificate-Authority certificates to enable HTTPS
RUN apk add --no-cache ca-certificates

ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .
RUN pwd
RUN ls -la

CMD ["/app/main"]
