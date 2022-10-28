package cli

import (
	"flag"
	"fmt"
	"log"

	"os"

	cluster "github.com/cclab.inu/testbed-mgmt/src/cluster"
)

// CLI responsible for processing command line arguments
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")

	fmt.Println("  create-cluster")
	fmt.Println("  delete-cluster")
	fmt.Println("  restart-cluster")

	fmt.Println("  deploy-pods")
	fmt.Println("  delete-pods")
	fmt.Println("  restart-pods")

	fmt.Println("  pull-image")
	fmt.Println("  delete-image")

	fmt.Println("  print-logs")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) createCluster() {
	cluster.CreateCluster()
}

func (cli *CLI) deleteCluster() {

}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	createCluster := flag.NewFlagSet("create-cluster", flag.ExitOnError)
	deleteCluster := flag.NewFlagSet("delete-cluster", flag.ExitOnError)

	switch os.Args[1] {
	case "create-cluster":
		err := createCluster.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "delete-cluster":
		err := deleteCluster.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if deleteCluster.Parsed() {
		cli.deleteCluster()
	}

	if createCluster.Parsed() {
		cli.createCluster()
	}
}
