package commands

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
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

// improved by Google Bard
func runCommand(command string, silent bool, workingDir string, logFile *os.File) error {
	var arguments []string
	if runtime.GOOS == "windows" {
		arguments = append(arguments, []string{"cmd", "/C"}...)
	} else {
		arguments = append(arguments, []string{"/bin/bash", "-c"}...)
	}
	arguments = append(arguments, command)

	cmd := exec.Command(arguments[0], arguments[1:]...)

	// Create a new process group for the subcommand
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

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

	// Capture system signals for better control
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		<-sigCh
		// Send SIGINT signal only to the subcommand process group
		syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
	}()

	err = cmd.Wait()
	close(sigCh)

	return err
}
