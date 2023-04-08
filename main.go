package main

import (
	"os"

	"github.com/STARRY-S/bilibili-danmaku-client/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(os.Args[1:]); err != nil {
		logrus.Fatal(err)
	}
}
