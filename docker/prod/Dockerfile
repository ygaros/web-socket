FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

RUN apk add sqlite

COPY ./sender ./

RUN go mod download

RUN go build -o application

COPY ./docker/ ./docker/

RUN ./docker/scripts/create-database.sh

EXPOSE 9876 8080

CMD [ "./application" ]  