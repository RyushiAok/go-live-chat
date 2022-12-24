FROM golang:1.19.4-alpine  

WORKDIR /go/src
COPY . .

EXPOSE 8080

RUN go mod download
RUN go build main.go