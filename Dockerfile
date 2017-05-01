FROM golang:1.7.5

RUN apt-get -y update
RUN apt-get -y install libxml2-dev

RUN mkdir /app

COPY . /go/src/app

RUN go get github.com/denisenkom/go-mssqldb
WORKDIR /go/src/app
RUN go build -o poller

CMD ["/go/src/app/poller"]
