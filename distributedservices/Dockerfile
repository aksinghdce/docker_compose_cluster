FROM golang:latest 
RUN mkdir /go/src/app 
ADD . /go/src/app
WORKDIR /go/src/app
EXPOSE 8080:8080
# RUN go build -o server distributedgrepserver/distributedgrepserver.go
# RUN go build -o client distributedgrep/distributedgrep.go
# CMD ["/go/src/app/server"]
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["clusterserver"]