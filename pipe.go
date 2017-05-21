package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"syscall"
	"regexp"
	"strconv"
)

type Pipe struct {
    Path  string
    Type  string
    Perms string
}

// Create and handle log data from pipe
func handlePipe(pipe Pipe) {
    var pipeOutput *os.File
    pipeExists := false

    // get pipe permissions
    if pipe.Perms == "" {
        pipe.Perms = "0600"
    }
    pipePerms, _ := strconv.ParseUint(pipe.Perms, 8, 32)

    // get pipe output destination
    switch pipe.Type {
    case "stdout":
        pipeOutput = os.Stdout
    case "stderr":
        pipeOutput = os.Stderr
    }

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
            fmt.Fprint(pipeOutput, message)
        }
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

			pipelist = append(pipelist, Pipe{Path:pipePath, Type:pipeType})
		} else {
			printHelp()
		}
	}

	return pipelist
}
