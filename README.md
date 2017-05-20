# go-logpipe

[![GitHub release](https://img.shields.io/github/release/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/releases)
[![license](https://img.shields.io/github/license/webdevops/go-replace.svg)](https://github.com/webdevops/go-replace/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-replace.svg?branch=master)](https://travis-ci.org/webdevops/go-replace)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-replace/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-replace/latest/total.svg)]()

Log collecting (using named pipes) utility written in golang

Inspired by https://github.com/abrander/logpipe

## Usage

```
Usage:
  go-logpipe [OPTIONS] Pipe...

Application Options:
      --permissions= Sets the permissions of the pipe (default: 0666)
  -V, --version      show version and exit
      --dumpversion  show only version number and exit

Help Options:
  -h, --help         Show this help message

Arguments:
  Pipe:              stdout:/path/to/pipe or stderr:/path/to/pipe
```
