FROM golang:1.23.5 AS builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client
RUN go build -o songs ./cmd/songs/main.go

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh 

CMD ["./songs"]