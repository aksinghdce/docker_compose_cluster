FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 3000:3000
RUN go build -o hello . 
CMD ["/app/hello"]