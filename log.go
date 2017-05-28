package main

import (
	"log"
)

type SyslogLogger struct {
	*log.Logger
}

func (SyslogLogger SyslogLogger) Verbose(message string) {
	if opts.Verbose {
		SyslogLogger.Println(message)
	}
}
