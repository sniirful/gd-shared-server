package colors

import "fmt"

func GreenBold(text string) string {
	return fmt.Sprintf("\x1b[32;1m%v\x1b[0m", text)
}

func RedBold(text string) string {
	return fmt.Sprintf("\x1b[31;1m%v\x1b[0m", text)
}

func Bold(text string) string {
	return fmt.Sprintf("\x1b[1m%v\x1b[0m", text)
}
