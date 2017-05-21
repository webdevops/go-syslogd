package main

import (
    "fmt"
    "os"
    "io/ioutil"
    "sync"
    "strings"
    flags "github.com/jessevdk/go-flags"
    yaml "gopkg.in/yaml.v2"
)

const (
    Name    = "go-syslogd"
    Author  = "webdevops.io"
    Version = "0.1.0"
)

var (
    argparser *flags.Parser
    configuration ConfigurationDefinition
)

type ConfigurationDefinition struct {
    Syslog Syslog
    Pipes []Pipe
}

var opts struct {
    Configuration           string   `short:"c"  long:"configuration"                 description:"Configuration file (yml)" default:"/etc/go-syslog.yml"`
    ShowVersion             bool     `short:"V"  long:"version"                       description:"show version and exit"`
    ShowOnlyVersion         bool     `           long:"dumpversion"                   description:"show only version number and exit"`
}

// handle special cli options
// eg. --help
//     --version
//     --path
//     --mode=...
func handleSpecialCliOptions(args []string) {
    // --dumpversion
    if (opts.ShowOnlyVersion) {
        fmt.Println(Version)
        os.Exit(0)
    }

    // --version
    if (opts.ShowVersion) {
        fmt.Println(fmt.Sprintf("%s version %s", Name, Version))
        fmt.Println(fmt.Sprintf("Copyright (C) 2017 %s", Author))
        os.Exit(0)
    }
}

// Parse configuration from yml file
func parseConfiguration() {
    confData, err := ioutil.ReadFile(opts.Configuration)
    if err != nil {
        panic(fmt.Sprintf("Failed opening configuration file %s: %v", opts.Configuration, err.Error()))
    }

    err = yaml.Unmarshal(confData, &configuration)
    if err != nil {
        panic(fmt.Sprintf("Unable to parse configuration file %s: %v", opts.Configuration, err.Error()))
    }

    if configuration.Syslog.Path != "" {
        configuration.Syslog.Filter.facility = 255

        // Facility filter
        for _, facility := range strings.Split(configuration.Syslog.Filter.Facility, ",") {
            if facilityId, ok := SyslogFacilityMap[facility]; ok {
                configuration.Syslog.Filter.facility -= facilityId
            }
        }
    }
}

func printMessage(dest *os.File, msg string) {
    fmt.Fprint(dest, msg)
}

// Prints help
func printHelp() {
    argparser.WriteHelp(os.Stdout)
    os.Exit(1)
}

// Main function
func main() {
    var wg sync.WaitGroup

    // init argument parser
    argparser = flags.NewParser(&opts, flags.Default)
    args, err := argparser.Parse()

    handleSpecialCliOptions(args)

    // check if there is an parse error
    if err != nil {
        if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
            os.Exit(0)
        } else {
            os.Exit(1)
        }
    }

    // parse yml configuration
    parseConfiguration()

    fmt.Println(fmt.Sprintf("Starting %s version %s", Name, Version))

    // init pipes
    for _, pipe := range configuration.Pipes {
        wg.Add(1)
        go func(pipe Pipe) {
            handlePipe(pipe)
            wg.Done()
        } (pipe);
    }

    // init syslog
    if configuration.Syslog.Path != "" {
        handleSyslog()
    }

    wg.Wait()
}
