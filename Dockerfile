FROM golang:1.13.5

# Enable Go modules
ENV GO111MODULE=on

ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o main .

CMD ["cd /app && ./main"]
