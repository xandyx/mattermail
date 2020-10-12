package main

import (
	"fmt"
	"os"

	"github.com/xandyx/mattermail/cmd"
)

func main() {
	if err := cmd.Execute(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
