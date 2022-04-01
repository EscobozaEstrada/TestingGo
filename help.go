package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func readOutput(rc io.ReadCloser, forward io.Writer, done func()) (wsURL string, _ error) {
	prefix := []byte("DevTools listening on")
	var accumulated bytes.Buffer
	bufr := bufio.NewReader(rc)
readLoop:
	for {
		line, err := bufr.ReadBytes('\n')
		if err != nil {
			return "", fmt.Errorf("chrome failed to start:\n%s",
				accumulated.Bytes())
		}
		if forward != nil {
			if _, err := forward.Write(line); err != nil {
				return "", err
			}
		}

		if bytes.HasPrefix(line, prefix) {
			line = line[len(prefix):]
			// use TrimSpace, to also remove \r on Windows
			line = bytes.TrimSpace(line)
			wsURL = string(line)
			break readLoop
		}
		accumulated.Write(line)
	}
	if forward == nil {
		// We don't need the process's output anymore.
		rc.Close()
	} else {
		// Copy the rest of the output in a separate goroutine, as we
		// need to return with the websocket URL.
		go func() {
			io.Copy(forward, bufr)
			done()
		}()
	}
	return wsURL, nil
}
