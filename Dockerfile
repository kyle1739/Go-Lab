FROM golang:alpine AS go-build 
WORKDIR /go/src/mml
ADD . .
RUN apk add --no-cache \
  git \ 
  && export GOPATH=/go/src \
  && go get github.com/gorilla/websocket \
  && go get github.com/olivere/elastic \
  && go get github.com/go-redis/redis \
  && go get github.com/go-sql-driver/mysql \
  && cd /go/src/mml/wankeliaoserver/server \
  && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mml 

FROM alpine:latest
COPY --from=go-build /go/src/mml/wankeliaoserver/ /  
WORKDIR /server
EXPOSE 80/tcp
ENTRYPOINT ["./mml"]
