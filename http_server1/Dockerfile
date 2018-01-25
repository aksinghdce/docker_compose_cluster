FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 8080:8080
RUN go build -o main main.go
RUN go build -o http_client http_client.go
CMD ["/app/main"]