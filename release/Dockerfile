FROM alpine

RUN apk --no-cache add --virtual .gosyslogd-deps \
        ca-certificates  \
        wget \
    && GOSYSLOGD_VERSION=0.1.0 \
    && wget -O /etc/go-syslog.yml https://raw.githubusercontent.com/webdevops/go-syslogd/master/etc/go-syslog.yml \
    && wget -O /usr/local/bin/go-syslogd https://github.com/webdevops/go-syslogd/releases/download/GOSYSLOGD_VERSION/go-syslogd-64-linux \
    && chmod +x /usr/local/bin/go-syslogd \
    && apk del .gosyslogd-deps

CMD ["go-syslogd"]
