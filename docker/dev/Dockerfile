FROM golang:1.18-alpine

WORKDIR /app

RUN apk add build-base

RUN apk add sqlite

COPY ./docker/ /app/docker

RUN chmod u+x ./docker/scripts/create-database.sh

RUN ./docker/scripts/create-database.sh

# ## code realoding - env
RUN go install github.com/cespare/reflex@latest

COPY ./docker/dev/reflex.conf /

COPY ./docker/dev/start.sh /

EXPOSE 9876 8080

ENTRYPOINT [ "reflex", "-c", "/reflex.conf" ]