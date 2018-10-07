package main

import (
	"github.com/fvdveen/mu2-encode/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.WithField("type", "main").Fatal(err)
	}
}
