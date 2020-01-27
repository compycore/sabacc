# Heroku doesn't seem to fully support multi-stage builds yet
FROM golang:1.13.5

# Enable Go modules
ENV GO111MODULE=on

ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .

CMD ["/app/main"]
