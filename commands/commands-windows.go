//go:build windows

package commands

import (
	"io"
	"os"
	"syscall"

	"github.com/ActiveState/termtest/conpty"
	"github.com/muesli/cancelreader"
	"golang.org/x/term"
)

// TODO: this function does not account for terminal
// resize at all, it only takes the first terminal
// size and uses it throughout the entire command
// execution
func runCommand(command string, workingDir string, logFile *os.File) error {
	// in stupid windows, the terminal size is gotten
	// from the stdout, not the stdin
	tw, th, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	// we do not want to close the pseudo-terminal here,
	// only when the process has ended, so it's handled
	// later in the code
	arguments := []string{`C:\Windows\System32\cmd.exe`, `/c`, command}
	pseudoTerminal, err := conpty.New(int16(tw), int16(th))
	if err != nil {
		return err
	}

	pid, _, err := pseudoTerminal.Spawn(
		arguments[0],
		arguments[1:],
		&syscall.ProcAttr{
			// we need the environment set, or
			// the path variable won't be set and
			// most commands won't work
			Env: os.Environ(),
			Dir: workingDir,
		},
	)
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// we need to wait for the process to finish in
	// a goroutine, so that we make sure that all
	// the stdout is copied synchronously and, most
	// importantly, before this function ends
	var commandError error
	go func() {
		_, commandError = process.Wait()
		pseudoTerminal.Close()
	}()

	// we need to use a cancelreader, because we need
	// to cancel any pending copy from stdin to cpty
	stdinReader, err := cancelreader.NewReader(os.Stdin)
	if err != nil {
		return err
	}
	defer stdinReader.Cancel()

	// io.Copy() for the stdin needs to be asynchronous
	// because we will cancel the reader at the end
	// of this function; it's more important that
	// the stdout is copied synchronously
	go io.Copy(pseudoTerminal.InPipe(), stdinReader)
	// we write the output of the pseudo-terminal to
	// both the stdout and the log file; we do not care
	// if the log file is nil, in that case the write
	// will just not happen
	io.Copy(io.MultiWriter(os.Stdout, logFile), pseudoTerminal.OutPipe())

	return commandError
}
