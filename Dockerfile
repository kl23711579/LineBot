FROM golang:alpine AS build-env
RUN apk add bash ca-certificates git gcc g++ libc-dev
RUN mkdir /app
ADD . /app
WORKDIR /app

# Force the go compiler to use module
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

#Final stage
FROM centurylink/ca-certs

ENTRYPOINT ["/app/app"]