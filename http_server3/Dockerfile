FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 8080:8080
RUN go build -o main . 
CMD ["/app/main"]