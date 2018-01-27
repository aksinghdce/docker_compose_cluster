FROM golang:latest 
RUN mkdir /go/src/app 
ADD . /go/src/app
WORKDIR /go/src/app
EXPOSE 8080:8080
RUN go build -o servergo server/server.go
RUN go build -o clientgo client/client.go
CMD ["/go/src/app/servergo"]