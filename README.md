### Prerequisites

- Go 1.18 or newer
- Java 8 on newer
- Node 16 or newer
- Docker 

## to run app server
standard version
```
docker build -t web-socket -f ./docker/Dockerfile .
docker run -it -p 9876:9876 -p 8080:8080 web-socket
```
with live code realoading
```
sudo docker build -t web-socket -f ./docker/dev/Dockerfile .
sudo docker run -v ~/training/go/notification/web-socket/sender:/app/sender -it -p 9876:9876 -p 8080:8080 web-socket
```
or
```
docker compose up --build
```
or (for non dokerized version)

```
go run ./sender/main.go
```
## to run app clients

```
go run ./go-client/main.go
```

```
node ./js-client/main.js
```

```
cd ./java-client
mvn package verify clean --fail-never
java -jar ./target/client-0.0.1-SNAPSHOT.jar
```
