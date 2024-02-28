package commands

import (
	"app/signals"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
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

func runCommand(command string, silent bool, workingDir string, logFile *os.File) error {
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

	// create a new process group for the subcommand so that
	// when Ctrl-C is pressed, it's not immediately passed
	// from father to child
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}

	onStopSignalling := signals.CaptureSIGINT(func() {
		syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
	})

	err := cmd.Wait()
	onStopSignalling()
	return err
}
