package pod

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/cclab.inu/testbed-mgmt/src/image"
)

// deploy-pods
func DeployPods() {
	var wg sync.WaitGroup
	wg.Add(1)

	image.PullImage()

	switch os.Args[2] {
	// Web Daemon
	case "nginx", "httpd", "mongo-express":
		// make image version
		image := os.Args[2] + ":"
		podName := os.Args[2] + "-"
		if len(os.Args) == 3 {
			image += "latest"
			podName += "latest"
		} else {
			image += os.Args[3]
			podName += os.Args[3]
		}

		// create deployment ; 1 pod
		runCmd := "kubectl create deployment " + podName + " --image=" + image + " --replicas=1 --port=80"
		cmd_pod := exec.Command("sh", "-c", runCmd)
		podOut, err := cmd_pod.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(podOut))

	// CMS
	case "joomla", "drupal", "wordpress":
		CMScmd := "helm install " + os.Args[2] + " bitnami/" + os.Args[2] + " --set service.port=8080"
		cmd_cms := exec.Command("sh", "-c", CMScmd)
		CMSout, err := cmd_cms.Output()
		if err != nil {
			log.Fatal(err)
		}
		output := string(CMSout)
		fmt.Println(output[:strings.Index(output, "**")])
		print("** " + os.Args[2] + " created ** \n")

	default:
		print("Check you Command\n")
	}
}

// delete-pods
func DeletePods() {
	var wg sync.WaitGroup
	wg.Add(1)

	switch os.Args[2] {
	case "all":
		cmd_deploy := exec.Command("sh", "-c", "kubectl delete deployment --all")
		cmd_deploy.Stdout = os.Stdout
		if err := cmd_deploy.Run(); err != nil {
			panic(err)
		}
		cmd_cms := exec.Command("sh", "-c", "helm uninstall $(helm ls --short)")
		cmd_cms.Stdout = os.Stdout
		if err := cmd_cms.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	// Web Daemon
	case "nginx", "httpd", "mongo":
		podName := os.Args[2] + "-"
		if len(os.Args) == 3 {
			podName += "latest"
		} else {
			podName += os.Args[3]
		}
		del_deploy := "kubectl delete deployment " + podName
		cmd_delete := exec.Command("sh", "-c", del_deploy)
		cmd_delete.Stdout = os.Stdout
		if err := cmd_delete.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	// CMS
	case "joomla", "drupal", "wordpress":
		del_release := "helm uninstall " + os.Args[2]
		cmd_delete := exec.Command("sh", "-c", del_release)
		cmd_delete.Stdout = os.Stdout
		if err := cmd_delete.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	default:
		fmt.Println("Check your Command")
	}
}

// restart-pods
func RestartPods() {
	var wg sync.WaitGroup
	wg.Add(1)

	switch os.Args[2] {
	case "nginx", "httpd", "mongo-express":
		podName := os.Args[2] + "-"
		if len(os.Args) == 3 {
			podName += "latest"
		} else {
			podName += os.Args[3]
		}
		re_deploy := "kubectl rollout restart deployment " + podName
		cmd_re := exec.Command("sh", "-c", re_deploy)
		cmd_re.Stdout = os.Stdout
		if err := cmd_re.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	case "joomla", "drupal", "wordpress":
		re_deploy := "kubectl rollout restart deployment " + os.Args[2]
		cmd_re := exec.Command("sh", "-c", re_deploy)
		cmd_re.Stdout = os.Stdout
		if err := cmd_re.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	default:
		fmt.Println("Check your Command")
	}
}
