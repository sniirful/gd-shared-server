package fileflags

import "os"

const (
	APPEND    = os.O_RDWR | os.O_CREATE | os.O_APPEND
	OVERWRITE = os.O_RDWR | os.O_CREATE | os.O_TRUNC
)
