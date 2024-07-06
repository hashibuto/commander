package commander

import (
	"fmt"
	"os"
	"strings"
)

var (
	C_RESET = "\x1b[0m"
	C_RED   = FgColor(255, 0, 0)
	C_GREEN = FgColor(0, 255, 0)
	C_BOLD  = "\x1b[1m"
)

func FgColor(red int, green int, blue int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", red, green, blue)
}

func Println(text ...string) {
	fmt.Printf("%s%s\n", strings.Join(text, ""), C_RESET)
}

func Errorln(text ...string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", C_RED, strings.Join(text, ""), C_RESET)
}

func Sprintf(text ...string) string {
	return fmt.Sprintf("%s%s", strings.Join(text, ""), C_RESET)
}
