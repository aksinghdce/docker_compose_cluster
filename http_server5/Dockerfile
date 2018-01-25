FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 8080:8080
WORKDIR /app
RUN go build -o clientgo client/main.go
RUN go build -o servergo server/main.go
CMD ["/app/servergo"]