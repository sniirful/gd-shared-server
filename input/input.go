package input

import "os"

func GetChar() string {
	var b []byte = make([]byte, 1)
	os.Stdin.Read(b)

	return string(b[0])
}
