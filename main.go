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
    Type  *os.File
}

var opts struct {
    Positional struct {
        Pipe []string `description:"stdout:/path/to/pipe or stderr:/path/to/pipe"`
    } `positional-args:"true" required:"yes"`

    PipePermissions         string   `           long:"permissions"                   description:"Sets the permissions of the pipe" default:"0666"`
    ShowVersion             bool     `short:"V"  long:"version"                       description:"show version and exit"`
    ShowOnlyVersion         bool     `           long:"dumpversion"                   description:"show only version number and exit"`
}

func createAndHandlePipe(pipe Pipe) {
    pipeExists := false
    pipePerms, _ := strconv.ParseUint(opts.PipePermissions, 10, 32)

    // check for existing file
    fileInfo, err := os.Stat(pipe.Path)

    if err == nil {
        if (fileInfo.Mode() & os.ModeNamedPipe) > 0 {
            pipeExists = true
        } else {
            fmt.Printf("%d != %d\n", os.ModeNamedPipe, fileInfo.Mode())
            panic(fmt.Sprintf("Pipe %s exists, but it's not a named pipe (FIFO)", pipe.Path))
        }
    }

    // Try to create pipe if needed
    if !pipeExists {
        err := syscall.Mkfifo(pipe.Path, uint32(pipePerms))
        if err != nil {
            panic(fmt.Sprintf("Creation of pipe %s failed: %v", pipe.Path, err.Error()))
        }
    }

    // Open pipe for reading
    fd, err := os.Open(pipe.Path)
    if err != nil {
        panic(fmt.Sprintf("Failed opening pipe %s: %v", pipe.Path, err.Error()))
    }
    defer fd.Close()
    reader := bufio.NewReader(fd)

    for {
        message, err := reader.ReadString(0xa)
        if err != nil && err != io.EOF {
            panic(fmt.Sprintf("Reading from pipe %s failed: %v", pipe.Path, err.Error()))
        }

        if message != "" {
            fmt.Fprint(pipe.Type, message)
        }
    }
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
        fmt.Println(fmt.Sprintf("go-logpipe version %s", Version))
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
            var pipeType *os.File
            m := pipeRegexp.FindStringSubmatch(line)

            pipePath := m[2]

            switch m[1] {
                case "stdout":
                    pipeType = os.Stdout
                case "stderr":
                    pipeType = os.Stderr
            }

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

    pipelist := buildPipelist(opts.Positional.Pipe)

    if (len(pipelist) == 0) {
        printHelp()
    }

    var wg sync.WaitGroup

    for _, pipe := range pipelist {
        wg.Add(1)
        go func(pipe Pipe) {
            createAndHandlePipe(pipe)
            wg.Done()
        } (pipe);
    }

    wg.Wait()
}
