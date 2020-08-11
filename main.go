package main

import (
	"github.com/fllaca/scheriff/cmd"
)

var (
	version = "development"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
