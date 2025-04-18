/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package streams

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// A stream processing app
type pipedApp struct {
	stdOutWriter io.Writer
}

func NewPipedApp(stdOutWriter io.Writer) App {
	return &pipedApp{
		stdOutWriter: stdOutWriter,
	}
}

func (a *pipedApp) Run() int {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			_, err := a.stdOutWriter.Write(scanner.Bytes())
			if err != nil {
				_, err := os.Stderr.Write([]byte(err.Error()))
				if err != nil {
					panic(fmt.Sprintf("Could not write to stderr: %s", err))
				}
				return 1
			}
			_, err = a.stdOutWriter.Write([]byte("\n"))
			if err != nil {
				_, err = os.Stderr.Write([]byte(err.Error()))
				if err != nil {
					panic(fmt.Sprintf("Could not write to stderr: %s", err))
				}
				return 1
			}
		}
		if err := scanner.Err(); err != nil {
			_, err = os.Stderr.Write([]byte(err.Error()))
			if err != nil {
				panic(fmt.Sprintf("Could not write to stderr: %s", err))
			}
			return 1
		}

	}
	return 0
}
