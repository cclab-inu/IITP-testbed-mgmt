package cli

import (
	"flag"
	"fmt"
	"log"

	"os"

	cluster "github.com/cclab.inu/testbed-mgmt/src/cluster"
	pods "github.com/cclab.inu/testbed-mgmt/src/pod"
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
	cluster.DeleteCluster()
}

func (cli *CLI) deployPod() {
	pods.DeployPods()
}

func (cli *CLI) deletePod() {
	pods.DeletePods()
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	createCluster := flag.NewFlagSet("create-cluster", flag.ExitOnError)
	deleteCluster := flag.NewFlagSet("delete-cluster", flag.ExitOnError)

	deployPod := flag.NewFlagSet("deploy-pods", flag.ExitOnError)
	deletePod := flag.NewFlagSet("deploy-pods", flag.ExitOnError)

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
	case "deploy-pods":
		err := deployPod.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "delete-pods":
		err := deletePod.Parse(os.Args[2:])
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

	if deployPod.Parsed() {
		cli.deployPod()
	}

	if deletePod.Parsed() {
		cli.deletePod()
	}
}
