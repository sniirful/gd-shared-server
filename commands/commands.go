package commands

import (
	"os"
)

func Run(command string) {
	runCommand(command, "", nil)
}

func RunWithWorkingDirAndLogFile(command, workingDir string, logFile *os.File) {
	runCommand(command, workingDir, logFile)
}
