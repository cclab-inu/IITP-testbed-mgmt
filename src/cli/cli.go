package cli

import (
	"flag"
	"fmt"
	"log"

	"os"

	cluster "github.com/cclab.inu/testbed-mgmt/src/cluster"
	logs "github.com/cclab.inu/testbed-mgmt/src/consumer"
	"github.com/cclab.inu/testbed-mgmt/src/image"
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

func (cli *CLI) restartCluster() {
	cluster.DeleteCluster()
	cluster.CreateCluster()
}

func (cli *CLI) deployPod() {
	pods.DeployPods()
}

func (cli *CLI) deletePod() {
	pods.DeletePods()
}

func (cli *CLI) restartPod() {
	pods.RestartPods()
}

func (cli *CLI) pullImage() {
	image.PullImage()
}

func (cli *CLI) deleteImage() {
	image.DeleteImage()
}

func (cli *CLI) printLogs() {
	logs.PrintLogs()
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	createCluster := flag.NewFlagSet("create-cluster", flag.ExitOnError)
	deleteCluster := flag.NewFlagSet("delete-cluster", flag.ExitOnError)
	restartCluster := flag.NewFlagSet("restart-cluster", flag.ExitOnError)

	deployPod := flag.NewFlagSet("deploy-pods", flag.ExitOnError)
	deletePod := flag.NewFlagSet("deploy-pods", flag.ExitOnError)
	restartPod := flag.NewFlagSet("restart-pods", flag.ExitOnError)

	pullImage := flag.NewFlagSet("pull-image", flag.ExitOnError)
	deleteImage := flag.NewFlagSet("delete-image", flag.ExitOnError)

	printLogs := flag.NewFlagSet("print-logs", flag.ExitOnError)

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
	case "restart-cluster":
		err := restartCluster.Parse(os.Args[2:])
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
	case "restart-pods":
		err := restartPod.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "pull-image":
		err := pullImage.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "delete-image":
		err := deleteImage.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print-logs":
		err := printLogs.Parse(os.Args[2:])
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

	if restartCluster.Parsed() {
		cli.restartCluster()
	}

	if deployPod.Parsed() {
		cli.deployPod()
	}

	if deletePod.Parsed() {
		cli.deletePod()
	}

	if restartPod.Parsed() {
		cli.restartPod()
	}

	if pullImage.Parsed() {
		cli.pullImage()
	}

	if deleteImage.Parsed() {
		cli.deleteImage()
	}

	if printLogs.Parsed() {
		cli.printLogs()
	}
}
