package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/geoffjay/plantd-core/core"
)

func main() {
	processArgs()
}

func processArgs() {
	if len(os.Args) > 1 {
		r := regexp.MustCompile("^-v$|(-{2})?version$")
		if r.Match([]byte(os.Args[1])) {
			fmt.Println(core.VERSION)
		}
		os.Exit(0)
	}
}
