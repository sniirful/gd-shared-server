package commands

import (
	"io"
	"os"
	"os/exec"
	"runtime"
)

func Run(command string) {
	runCommand(command, false, "", nil)
}

func RunSilent(command string) {
	runCommand(command, true, "", nil)
}

func RunWithWorkingDirAndLogFile(command, workingDir string, logFile *os.File) {
	runCommand(command, false, workingDir, logFile)
}

func runCommand(command string, silent bool, workingDir string, logFile *os.File) {
	var arguments []string
	if runtime.GOOS == "windows" {
		arguments = append(arguments, []string{"cmd", "/C"}...)
	} else {
		arguments = append(arguments, []string{"/bin/bash", "-c"}...)
	}
	arguments = append(arguments, command)

	cmd := exec.Command(arguments[0], arguments[1:]...)

	writers := []io.Writer{}
	if !silent {
		cmd.Stdin = os.Stdin

		writers = append(writers, os.Stdout)
	}
	if logFile != nil {
		writers = append(writers, logFile)
	}
	writer := io.MultiWriter(writers...)
	cmd.Stdout = writer
	cmd.Stderr = writer
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	cmd.Run()
}
