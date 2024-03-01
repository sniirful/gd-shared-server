//go:build windows

package commands

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func runCommand(command string, workingDir string, logFile *os.File) error {
	// TODO: add support for Windows
	var arguments []string
	if runtime.GOOS == "windows" {
		arguments = append(arguments, []string{"cmd", "/C"}...)
	} else {
		arguments = append(arguments, []string{"/bin/bash", "-c"}...)
	}
	arguments = append(arguments, command)

	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stdin = os.Stdin

	writers := []io.Writer{}
	writers = append(writers, os.Stdout)
	if logFile != nil {
		writers = append(writers, logFile)

		// we make sure to copy the command stdin to
		// the logfile as well, to make it more
		// readable
		go io.Copy(logFile, cmd.Stdin)
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
