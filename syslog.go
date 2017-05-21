package main

import (
    "fmt"
    "os"
    syslog "gopkg.in/mcuadros/go-syslog.v2"
)


type Syslog struct {
    Path string
    Filter struct {
        Facility string
        facility int
    }
    Output struct {
        Template string
    }
}


var SyslogFacilityMap = map[string]int {
    "kern": 0,
    "user": 1,
    "mail": 2,
    "daemon": 3,
    "auth": 4,
    "syslog": 5,
    "lpr": 6,
    "news": 7,
    "uucp": 8,
    "cron": 9,
    "authpriv": 10,
    "ftp": 11,
    "ntp": 12,
    "security": 13,
    "console": 14,
    "solaris-cron": 15,
    "local0": 16,
    "local1": 17,
    "local2": 18,
    "local3": 19,
    "local4": 20,
    "local5": 21,
    "local6": 22,
    "local7": 23,
}


func handleSyslog() {
    LoggerStdout.Verbose(fmt.Sprintf(" -> starting syslog daemon (%s)", configuration.Syslog.Path))

    // Check if syslog path exists, remove if already existing
    _, err := os.Stat(configuration.Syslog.Path)
    if err == nil {
        os.Remove(configuration.Syslog.Path)
    }

    channel := make(syslog.LogPartsChannel)
    handler := syslog.NewChannelHandler(channel)

    server := syslog.NewServer()
    server.SetFormat(syslog.Automatic)
    server.SetHandler(handler)
    server.ListenUnixgram(configuration.Syslog.Path)
    server.Boot()

    go func(channel syslog.LogPartsChannel) {
        for logParts := range channel {
            facilityId := uint(logParts["facility"].(int))

            // facility filter
            if hasBit(configuration.Syslog.Filter.facility, facilityId) == false {
                continue
            }

            //fmt.Println(logParts)

            // build message
            message := fmt.Sprintf("%s %s", logParts["hostname"], logParts["content"])

            // custom template
            if configuration.Syslog.Output.Template != "" {
                message = fmt.Sprintf(configuration.Syslog.Output.Template, message)
            }

            LoggerStdout.Println(message)
        }
    }(channel)

    server.Wait()
}
