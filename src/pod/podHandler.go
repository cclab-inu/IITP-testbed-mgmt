package pod

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/rs/zerolog/log"
)

func runCMD(cmd *exec.Cmd) {
	out, err := cmd.Output()
	if err != nil {
		log.Err(err)
		return
	}
	log.Info().Msg(string(out))
}

// deploy-pods
func DeployPods() {
	log.Info().Msg("Deploying pod: " + os.Args[2])

	var wg sync.WaitGroup
	wg.Add(1)

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
		cmd := "kubectl create deployment " + podName + " --image=" + image + " --replicas=1 --port=80"
		runCMD(exec.Command("sh", "-c", cmd))

	// CMS
	case "joomla", "drupal", "wordpress":
		cmd := "helm install " + os.Args[2] + " bitnami/" + os.Args[2] + " --set service.port=8080"
		runCMD(exec.Command("sh", "-c", cmd))

	default:
		print("Check your Command\n")
	}
}

// delete-pods
func DeletePods() {
	log.Info().Msg("Deleting pod: " + os.Args[2])

	var wg sync.WaitGroup
	wg.Add(1)

	switch os.Args[2] {
	case "all":
		cmd_deploy := exec.Command("sh", "-c", "kubectl delete deployment --all")
		runCMD(cmd_deploy)

		cmd_cms := exec.Command("sh", "-c", "helm uninstall $(helm ls --short)")
		runCMD(cmd_cms)

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
		runCMD(cmd_delete)

		wg.Done()
	// CMS
	case "joomla", "drupal", "wordpress":
		del_release := "helm uninstall " + os.Args[2]
		cmd_delete := exec.Command("sh", "-c", del_release)
		runCMD(cmd_delete)

		wg.Done()
	default:
		fmt.Println("Check your Command")
	}
}

// restart-pods
func RestartPods() {
	log.Info().Msg("Restarting pod: " + os.Args[2])

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
		runCMD(cmd_re)

		wg.Done()
	case "joomla", "drupal", "wordpress":
		re_deploy := "kubectl rollout restart deployment " + os.Args[2]
		cmd_re := exec.Command("sh", "-c", re_deploy)
		runCMD(cmd_re)

		wg.Done()
	default:
		fmt.Println("Check your Command")
	}
}
