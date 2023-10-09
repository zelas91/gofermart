FROM golang:1.20
LABEL authors="zelas"
COPY ./ ./
ENV GOPATH=/
RUN go mod download

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install migrate
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.12.1/wait ./wait
RUN chmod +x ./wait



