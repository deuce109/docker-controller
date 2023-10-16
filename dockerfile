FROM golang:alpine as build

WORKDIR /build

COPY ./go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go get github.com/deuce109/docker-controller/v2/handlers

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-controller


CMD [ "/docker-controller" ]