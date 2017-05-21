# go-syslogd

[![GitHub release](https://img.shields.io/github/release/webdevops/go-syslogd.svg)](https://github.com/webdevops/go-syslogd/releases)
[![license](https://img.shields.io/github/license/webdevops/go-syslogd.svg)](https://github.com/webdevops/go-syslogd/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-syslogd.svg?branch=master)](https://travis-ci.org/webdevops/go-syslogd)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-syslogd/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-syslogd/latest/total.svg)]()

Syslog daemon written in golang which also provides named pipes (FIFO)

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
