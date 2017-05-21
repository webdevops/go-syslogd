# go-syslogd

[![GitHub release](https://img.shields.io/github/release/webdevops/go-syslogd.svg)](https://github.com/webdevops/go-syslogd/releases)
[![license](https://img.shields.io/github/license/webdevops/go-syslogd.svg)](https://github.com/webdevops/go-syslogd/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-syslogd.svg?branch=master)](https://travis-ci.org/webdevops/go-syslogd)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-syslogd/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-syslogd/latest/total.svg)]()

Syslog daemon written in golang which also provides named pipes (FIFO)

This daemon will collect written logs and syslog messages and writes them to STDOUT and STDERR (eg. for usage in Docker)

Inspired by https://github.com/abrander/logpipe

## Usage

```
Usage:
  go-syslogd [OPTIONS]

Application Options:
  -c, --configuration= Configuration file (yml) (default: /etc/go-syslog.yml)
  -V, --version        show version and exit
      --dumpversion    show only version number and exit

Help Options:
  -h, --help           Show this help message

```

## Configuration

see [etc/go-syslog.yml](etc/go-syslog.yml) for an example.


## Installation

```bash
GOSYSLOGD_VERSION=0.1.0 \
&& wget -O /etc/go-syslog.yml https://raw.githubusercontent.com/webdevops/go-syslogd/master/etc/go-syslog.yml \
&& wget -O /usr/local/bin/go-syslogd https://github.com/webdevops/go-syslogd/releases/download/GOSYSLOGD_VERSION/go-syslogd-64-linux \
&& chmod +x /usr/local/bin/go-syslogd
```

## Docker images

| Image                         | Description                                                         |
|:------------------------------|:--------------------------------------------------------------------|
| `webdevops/go-syslogd:latest` | Latest release, binary only                                         |
| `webdevops/go-syslogd:master` | Current development version in branch `master`, with golang runtime |
