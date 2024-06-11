package size

import (
	"strings"

	"github.com/inhies/go-bytesize"
)

var sizesMap = map[string]string{
	"KB": "kiB",
	"MB": "MiB",
	"GB": "GiB",
	"TB": "TiB",
}

// this implementation is necessary to give the user
// accurate sizes; the bytesize library uses "MB" for "MiB"
func Parse(size int64) string {
	res := bytesize.New(float64(size)).String()
	for k, v := range sizesMap {
		if strings.Contains(res, k) {
			return strings.Replace(res, k, v, 1)
		}
	}
	return res
}
