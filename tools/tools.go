package tools

import (
	"strings"
)

func StringToArgs(str string) (args []string) {
	args = make([]string, 0)

	for _, s := range strings.Split(str, " ") {
		if len(s) > 0 {
			append(args, s)
		}
	}

	return
}