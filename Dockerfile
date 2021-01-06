FROM golang:alpine

WORKDIR /app
ADD . /app/
RUN go build -o /usr/bin/lookupbot main.go
RUN chmod +x /usr/bin/lookupbot

CMD ["/usr/bin/lookupbot"]

FROM golang:alpine
