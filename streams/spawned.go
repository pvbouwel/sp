/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package streams

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

//A stream processing app
type spawnedApp struct {
	stdOutWriter io.Writer
	stdErrWriter io.Writer

	appName string
	appArgs []string
}

func NewSpawnedApp(stdOutWriter io.Writer, StdErrWriter io.Writer, appName string, AppArgs []string) App{
	return &spawnedApp{
		stdOutWriter: stdOutWriter,
		stdErrWriter: StdErrWriter,
		appName: appName,
		appArgs: AppArgs,
	}
}


func (a *spawnedApp) Run() int {
	appPath, err := exec.LookPath(a.appName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not locate app %s: %s", a.appName, err)
		return 1
	}
	prog := exec.Command(appPath, a.appArgs...)
	prog.Stdout = a.stdOutWriter
	prog.Stderr = a.stdErrWriter
	err = prog.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Spawned app %s got error: %s", appPath, err)
		return 1
	}
	return 0
}
