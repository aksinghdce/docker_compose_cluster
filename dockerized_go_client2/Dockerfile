FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o hello_from_client . 
CMD ["/app/hello_from_client"]