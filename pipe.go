package main

import (
	"fmt"
	"os"
	"io"
	"bufio"
	"syscall"
	"strconv"
)

type Pipe struct {
    Path  string
    Type  string
    Perms string

    Output struct {
        Template string
    }
}

// Create and handle log data from pipe
func handlePipe(pipe Pipe) {
    pipeExists := false

    LoggerStdout.Verbose(fmt.Sprintf(" -> starting named pipe (%s)", pipe.Path))

    // get pipe permissions
    if pipe.Perms == "" {
        pipe.Perms = "0600"
    }
    pipePerms, _ := strconv.ParseUint(pipe.Perms, 8, 32)

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
            panic(fmt.Sprintf("Creation of named pipe %s failed: %v", pipe.Path, err.Error()))
        }
    }

    // Open pipe for reading
    fd, err := os.Open(pipe.Path)
    if err != nil {
        panic(fmt.Sprintf("Failed opening named pipe %s: %v", pipe.Path, err.Error()))
    }
    defer fd.Close()
    reader := bufio.NewReader(fd)

    // loop messages
    for {
        message, err := reader.ReadString(0xa)
        if err != nil && err != io.EOF {
            panic(fmt.Sprintf("Reading from named pipe %s failed: %v", pipe.Path, err.Error()))
        }

        if message != "" {

            if pipe.Output.Template != "" {
                message = fmt.Sprintf(pipe.Output.Template, message)
            }

            switch pipe.Type {
                case "stdout":
                    LoggerStdout.Print(message)
                case "stderr":
                    LoggerStderr.Print(message)
            }
        }
    }
}
