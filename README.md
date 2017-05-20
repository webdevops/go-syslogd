# go-logd

[![GitHub release](https://img.shields.io/github/release/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/releases)
[![license](https://img.shields.io/github/license/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-replace.svg?branch=master)](https://travis-ci.org/webdevops/go-replace)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-replace/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-replace/latest/total.svg)]()

Log daemon written in golang which provides syslog and named pipes (FIFO)

Inspired by https://github.com/abrander/logpipe

## Usage

```
Usage:
  go-logd [OPTIONS]

Application Options:
      --syslog       Provide syslog server
      --pipe=        Setup file based named pipe for collecting log informations (eg. stdout:/path/to/file or stderr:/path/to/file)
      --permissions= Sets the permissions of the pipe (default: 0666)
  -V, --version      show version and exit
      --dumpversion  show only version number and exit

Help Options:
  -h, --help         Show this help message
```
