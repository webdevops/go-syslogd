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
}

var SyslogFacilityLookup = map[int]string {
    0: "kern",
    1: "user",
    2: "mail",
    3: "daemon",
    4: "auth",
    5: "syslog",
    6: "lpr",
    7: "news",
    8: "uucp",
    9: "cron",
    10: "authpriv",
    11: "ftp",
    12: "ntp",
    13: "security",
    14: "console",
}

func handleSyslog() {
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
            // facility filter
            if configuration.Syslog.Filter.facility & logParts["facility"].(int) == 0 {
                continue
            }

            facilityId := logParts["facility"].(int)
            facility := "custom"
            if val, ok := SyslogFacilityLookup[facilityId]; ok {
                facility = val
            }

            message := fmt.Sprintf("%s: %s", facility, logParts["content"])
            fmt.Println(message)
        }
    }(channel)

    server.Wait()
}
