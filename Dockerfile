FROM golang:1.10.0

ADD . /go/src/sevki.org/cf-controller

RUN go get sevki.org/cf-controller

CMD cf-controller
