package image

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// pull-image
func PullImage() {
	var wg sync.WaitGroup
	wg.Add(1)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

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

		// image pull
		pullCmd := "docker pull " + image
		cmd_img := exec.Command("sh", "-c", pullCmd)
		cmd_img.Stderr = os.Stderr
		pullOut, err := cmd_img.Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(pullOut))

		// image save
		saveCmd := "docker save -o " + dir + "/template/" + podName + ".tar " + image
		cmd_saveImg := exec.Command("sh", "-c", saveCmd)
		_, err = cmd_saveImg.Output()
		if err != nil {
			panic(err)
		}
		print("Successfully saved Image: " + podName + ".tar \n")

	// CMS
	case "joomla", "drupal", "wordpress":
		CMScmd := "helm pull bitnami/" + os.Args[2] + " -d " + dir + "/template/"
		cmd_chart := exec.Command("sh", "-c", CMScmd)
		_, err = cmd_chart.Output()
		if err != nil {
			panic(err)
		}
		print("Successfully saved Charts: " + os.Args[2] + "\n")

	default:
		print("Check you Command\n")
	}
}

func DeleteImage() {
	var wg sync.WaitGroup
	wg.Add(1)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	switch os.Args[2] {
	case "all":
		allCmd := "fm " + dir + "/template/*"
		cmd_delall := exec.Command("sh", "-c", allCmd)
		_, err := cmd_delall.Output()
		if err != nil {
			panic(err)
		}

	// Web Daemon
	case "nginx", "httpd", "mongo-express":
		// make image version
		image := os.Args[2]
		fileName := os.Args[2] + "-"
		if len(os.Args) == 3 {
			image += ":latest"
			fileName += "latest.tar"
		} else {
			if os.Args[3] == "all" {
				if os.Args[3] == "all" {
					fileCmd := "find " + dir + "/template/" + " -type f -name " + "\"*" + os.Args[2] + "*\"" + " -exec rm {} \\;"
					cmd_delfile := exec.Command("sh", "-c", fileCmd)
					_, err := cmd_delfile.Output()
					if err != nil {
						panic(err)
					}

					docCmd := "docker rmi `docker images | awk '$1 ~ /" + os.Args[2] + "/ {print $3}'`"
					cmd_deldoc := exec.Command("sh", "-c", docCmd)
					_, err = cmd_deldoc.Output()
					if err != nil {
						panic(err)
					}
				}
				print("Successfully deleted all of Images: " + os.Args[2] + " \n")
				os.Exit(1)
			}
			image += ":" + os.Args[3]
			fileName += os.Args[3]
		}
		// docker image delete
		rmiCmd := "docker rmi " + image
		cmd_delImg := exec.Command("sh", "-c", rmiCmd)
		_, err := cmd_delImg.Output()
		if err != nil {
			panic(err)
		}
		// template file delete
		selCmd := "find " + dir + "/template/" + " -type f -name " + fileName + " -exec rm {} \\;"
		cmd_delSel := exec.Command("sh", "-c", selCmd)
		_, err = cmd_delSel.Output()
		if err != nil {
			panic(err)
		}
		print("Successfully deleted Image: " + fileName + " \n")
		wg.Done()

	// CMS
	case "joomla", "drupal", "wordpress":
		fileCmd := "find " + dir + "/template/" + " -type f -name " + "\"*" + os.Args[2] + "*\"" + " -exec rm {} \\;"
		cmd_delfile := exec.Command("sh", "-c", fileCmd)
		_, err := cmd_delfile.Output()
		if err != nil {
			panic(err)
		}
		print("Successfully deleted Image: " + os.Args[2] + " \n")
		wg.Done()
	default:
		fmt.Println("Check your Command")
	}

}
