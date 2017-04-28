FROM golang:1.7.5

RUN apt-get -y update
RUN apt-get -y install libxml2-dev

RUN mkdir /app

RUN mkdir /opt/clidriver

COPY ./clidriver /opt/clidriver
COPY ./clidriver/lib/*.so* /usr/lib/
COPY ./clidriver/include/* /usr/include/
COPY . /go/src/app

RUN cd /opt/clidriver

RUN go get bitbucket.org/phiggins/db2cli
WORKDIR /go/src/app
RUN go build -o poller
ENV GODEBUG=cgocheck=0

ENV curDir=/opt/clidriver
ENV LD_LIBRARY_PATH="/opt/clidriver/lib":"$LD_LIBRARY_PATH"
ENV IBM_DB_DIR="$curDir"
ENV IBM_DB_HOME="$curDir"
ENV IBM_DB_LIB="$curDir/lib"
ENV IBM_DB_INCLUDE="$curDir/include"
ENV DB2_HOME="$curDir/include"
ENV INCLUDE="$curDir/include"
ENV DB2LIB="$curDir/lib"
CMD ["/go/src/app/poller"]
