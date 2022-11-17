package main

import (
	"os"

	"github.com/cclab.inu/testbed-mgmt/src/cli"
)

func init() {
	os.Setenv("HOME", "/home/cclab")
	os.Setenv("USER", "cclab")
}

func main() {
	cli := cli.CLI{}
	cli.Run()
}
