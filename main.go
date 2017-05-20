package main

import (
    "fmt"
    "os"
    "io"
    "bufio"
    "syscall"
    "sync"
    "regexp"
    "strconv"
    flags "github.com/jessevdk/go-flags"
    syslog "gopkg.in/mcuadros/go-syslog.v2"
)

const (
    Author  = "webdevops.io"
    Version = "1.0.0"
    PipeArgumentRegexp = "^(stdout|stderr):(.*)$"
)

var (
    argparser *flags.Parser
)

type Pipe struct {
    Path  string
    Type  string
}

var opts struct {
    Syslog                  bool     `           long:"syslog"                        description:"Provide syslog server"`
    Pipes                 []string   `           long:"pipe"                          description:"Setup file based named pipe for collecting log informations (eg. stdout:/path/to/file or stderr:/path/to/file)"`
    PipePermissions         string   `           long:"permissions"                   description:"Sets the permissions of the pipe" default:"0666"`
    ShowVersion             bool     `short:"V"  long:"version"                       description:"show version and exit"`
    ShowOnlyVersion         bool     `           long:"dumpversion"                   description:"show only version number and exit"`
}

func createAndHandlePipe(path string, output *os.File) {
    pipeExists := false
    pipePerms, _ := strconv.ParseUint(opts.PipePermissions, 10, 32)

    // check for existing file
    fileInfo, err := os.Stat(path)

    if err == nil {
        if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
            pipeExists = true
        } else {
            fmt.Printf("%d != %d\n", os.ModeNamedPipe, fileInfo.Mode())
            panic(fmt.Sprintf("Pipe %s exists, but it's not a named pipe (FIFO)", path))
        }
    }

    // Try to create pipe if needed
    if !pipeExists {
        err := syscall.Mkfifo(path, uint32(pipePerms))
        if err != nil {
            panic(fmt.Sprintf("Creation of pipe %s failed: %v", path, err.Error()))
        }
    }

    // Open pipe for reading
    fd, err := os.Open(path)
    if err != nil {
        panic(fmt.Sprintf("Failed opening pipe %s: %v", path, err.Error()))
    }
    defer fd.Close()
    reader := bufio.NewReader(fd)

    for {
        message, err := reader.ReadString(0xa)
        if err != nil && err != io.EOF {
            panic(fmt.Sprintf("Reading from pipe %s failed: %v", path, err.Error()))
        }

        if message != "" {
            fmt.Fprint(output, message)
        }
    }
}

func createAndHandleSyslog(path string) {
    // Check if syslog path exists, remove if already existing
    _, err := os.Stat(path)
    if err == nil {
        os.Remove(path)
    }

    channel := make(syslog.LogPartsChannel)
    handler := syslog.NewChannelHandler(channel)

    server := syslog.NewServer()
    server.SetFormat(syslog.Automatic)
    server.SetHandler(handler)
    server.ListenUnixgram(path)
    server.Boot()

    go func(channel syslog.LogPartsChannel) {
        for logParts := range channel {
            message := fmt.Sprintf("%s", logParts["content"])
            fmt.Println(message)
        }
    }(channel)

    server.Wait()
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
        fmt.Println(fmt.Sprintf("go-logd version %s", Version))
        fmt.Println(fmt.Sprintf("Copyright (C) 2017 %s", Author))
        os.Exit(0)
    }
}

// Build a pipe array list from command line arguments
func buildPipelist(args []string) ([]Pipe) {
    var pipelist []Pipe
    pipeRegexp := regexp.MustCompile(PipeArgumentRegexp)

    for _, line := range args {
        // check if line is matching our regexp
        if pipeRegexp.MatchString(line) == true {
            m := pipeRegexp.FindStringSubmatch(line)

            pipeType := m[1]
            pipePath := m[2]

            pipelist = append(pipelist, Pipe{pipePath, pipeType})
        } else {
            printHelp()
        }
    }

    return pipelist
}

// Prints hel
func printHelp() {
    argparser.WriteHelp(os.Stdout)
    os.Exit(1)
}

func main() {
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

    pipelist := buildPipelist(opts.Pipes)

    if (len(pipelist) == 0 && opts.Syslog == false) {
        printHelp()
    }

    var wg sync.WaitGroup

    for _, pipe := range pipelist {
        wg.Add(1)
        go func(pipe Pipe) {
            switch pipe.Type {
                case "stdout":
                    createAndHandlePipe(pipe.Path, os.Stdout)
                case "stderr":
                    createAndHandlePipe(pipe.Path, os.Stdout)
            }
            wg.Done()
        } (pipe);
    }

    if opts.Syslog {
        createAndHandleSyslog("/dev/log")
    }

    wg.Wait()
}
