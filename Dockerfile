FROM golang:1.18.4

ENV GOPACKAGE github.com/Barzahlen/nagios-dnsblklist

ADD . /go/src/$GOPACKAGE

WORKDIR /go/src/$GOPACKAGE

RUN go install github.com/mitchellh/gox@latest
