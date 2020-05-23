FROM golang:1.13-alpine3.10

RUN apk update
RUN apk add ca-certificates curl bash git
RUN apk add g++ gcc pkgconf make zlib-dev

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure --reconfigure --prefix /usr
RUN make
RUN make install

WORKDIR /go/src/kafka-producer

# cache dependencies
ADD go.* ./
RUN go mod download

ADD . ./
RUN go build -o kafka-producer cmd/kafka-producer/*
RUN mv kafka-producer /bin/kafka-producer

RUN go version

ENTRYPOINT [ "/bin/kafka-producer" ]
