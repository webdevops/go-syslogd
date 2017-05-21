FROM golang:alpine

COPY . /go/src/go-syslogd

WORKDIR /go/src/go-syslogd

RUN apk --no-cache add --virtual .gosyslogd-deps git \
    && go get \
    && go build \
    && mv go-syslogd /usr/local/bin \
    && rm -rf /go/src/ \
    && apk del .gosyslogd-deps

CMD ["go-syslogd"]
