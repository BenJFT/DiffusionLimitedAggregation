package util

import (
	"os"
	"bufio"
	"strings"
)

func StringToArgs(str string) (args []string) {
	args = make([]string, 0)

	for _, s := range strings.Split(str, " ") {
		if len(s) > 0 {
			args = append(args, s)
		}
	}

	return
}

var scanner = bufio.NewScanner(os.Stdin)
func ReadStrOrEmpty() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}