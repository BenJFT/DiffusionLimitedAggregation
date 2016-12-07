package tools

import (
	"strings"
	"bufio"
	"fmt"
	"os"
)

var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)

func SingleSpace(str string) string  {

	for tmpStr := strings.Replace(str, "  ", " ", -1); tmpStr != str; tmpStr = strings.Replace(str, "  ", " ", -1) {
		str = tmpStr
	}

	return str
}


func ReadInt(a *int64) (err error) {
	_, err = fmt.Scanf("%d\n", a)
	return
}

func ReadStrOrEmpty() string {
	b := scanner.Scan()
	if b {
		return scanner.Text()
	} else {
		return ""
	}
}