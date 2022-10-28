package main

import (
	"os"

	CLI "github.com/cclab.inu/testbed-mgmt/src/cli"
)

func init() {
	os.Setenv("HOME", "/home/cclab")
	os.Setenv("USER", "cclab")
}

func main() {
	cli := CLI.CLI{}
	cli.Run()
}
