//go:build !windows

package commands

import (
	"app/signals"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

func runCommand(command string, workingDir string, logFile *os.File) error {
	arguments := []string{"/bin/bash", "-c", command}
	cmd := exec.Command(arguments[0], arguments[1:]...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	pseudoTerminal, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	// we need to have a new goroutine here because
	// we want to close the pty only after the command
	// has finished running, but at the same time we
	// need to make sure that all the output is copied
	// to the tty synchronously and before closing
	// this function; of waiting in this  function and
	// calling io.Copy() in another goroutine makes
	// the command write to stdout sometimes even after
	// the function has returned, so we do things the
	// other way around
	var commandError error
	go func() {
		commandError = cmd.Wait()
		pseudoTerminal.Close()
	}()

	// this block of code makes sure the
	// pseudo-terminal is the same size as the
	// actual terminal
	onStopSignalling := signals.CaptureSIGWINCH(func() {
		pty.InheritSize(os.Stdin, pseudoTerminal)
	})
	defer onStopSignalling()

	// this sets stdin in raw mode, meaning that
	// commands like Ctrl-C will be captured
	// by stdin and later sent to the running
	// command
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// we need to use a cancelreader, because we need
	// to cancel any pending copy from stdin to pty
	stdinReader, err := cancelreader.NewReader(os.Stdin)
	if err != nil {
		return err
	}
	defer stdinReader.Cancel()

	// io.Copy() for the stdin needs to be asynchronous
	// because we will cancel the reader at the end
	// of this function; it's more important that
	// the stdout is copied synchronously
	go io.Copy(pseudoTerminal, stdinReader)
	// we write the output of the pseudo-terminal to
	// both the stdout and the log file; we do not care
	// if the log file is nil, in that case the write
	// will just not happen
	io.Copy(io.MultiWriter(os.Stdout, logFile), pseudoTerminal)

	return commandError
}
