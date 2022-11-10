package pod

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// Deployment Yaml Struct
type Yaml2Go struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

// Metadata
type Metadata struct {
	Name   string `yaml:"name"`
	Labels Labels
}

type Spec struct {
	Replicas int
	Selector Selector
	Template Template
}

type Selector struct { //spec:Selector
	MatchLabels Labels `yaml:"matchLabels"`
}

// Labels
type Labels struct {
	App string `yaml:"app"`
}

// spec:Template
type Template struct {
	Metadata TemMetadata `yaml:"metadata"`
	Spec     TemSpec     `yaml:"spec"`
}

type TemMetadata struct {
	Name   string
	Labels Labels
}

type TemSpec struct {
	Containers []Containers
}

// Containers
type Containers struct {
	Name  string  `yaml:"name"`
	Image string  `yaml:"image"`
	Ports []Ports `yaml:"ports"`
}

// Ports
type Ports struct {
	ContainerPort int `yaml:"containerPort"`
}

func deployPods() {
	var wg sync.WaitGroup
	wg.Add(2)

	var imgVER string
	fmt.Print("Image Version : ")

	dpReader := bufio.NewReader(os.Stdin)
	imgVER, err := dpReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	imgVER = strings.TrimSpace(imgVER)

	slice := strings.Split(imgVER, " ")
	imgVER = "docker pull " + slice[0] + ":" + slice[1]

	cmd := exec.Command("sh", "-c", imgVER)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Deployment Yaml File
	pod := Yaml2Go{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: Metadata{
			Name: slice[0] + "-" + slice[1] + "-deploy",
			Labels: Labels{
				App: "testbed-" + slice[0],
			}},
		Spec: Spec{
			Replicas: 3,
			Selector: Selector{
				MatchLabels: Labels{
					App: "testbed-" + slice[0],
				}},
			Template: Template{
				Metadata: TemMetadata{
					Name: "testbed-" + slice[0] + "-pod",
					Labels: Labels{
						App: "testbed-" + slice[0],
					}},
				Spec: TemSpec{
					Containers: []Containers{
						Containers{
							Name:  slice[0],
							Image: slice[0] + ":" + slice[1],
							Ports: []Ports{
								Ports{
									ContainerPort: 80,
								}}}}}}}}

	yamlData, err := yaml.Marshal(&pod)
	if err != nil {
		fmt.Println("Error while Marshaling. %v", err)
	}

	fileName := "deployment-" + slice[0] + ".yaml"
	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	println()
	deploymentPod := "kubectl apply -f " + fileName
	cmd_pod := exec.Command("sh", "-c", deploymentPod)
	cmd_pod.Stdout = os.Stdout
	if err := cmd_pod.Run(); err != nil {
		panic(err)
	}
}

func deletePods() {
	var wg sync.WaitGroup
	wg.Add(1)

	//deployPods()
	fmt.Print("Type(all/select) : ")

	dtReader := bufio.NewReader(os.Stdin)
	types, err := dtReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	types = strings.TrimSpace(types)

	switch types {
	case "all":
		cmd_all := exec.Command("sh", "-c", "kubectl delete deployment --all")
		cmd_all.Stdout = os.Stdout
		if err := cmd_all.Run(); err != nil {
			panic(err)
		}
		wg.Done()
	case "select":
		cmd_get, _ := exec.Command("sh", "-c", "kubectl get deployment -o name").Output()

		var deployName string
		deployCheck := strings.Split(string(cmd_get), "deployment.apps/")
		for _, str := range deployCheck {
			deployName = strings.Trim(str, " ")
			fmt.Print(deployName)
		}
		if len(deployName) <= 0 {
			println("check")
			wg.Done()
			break
		} else {
			var imgVER string
			println()
			fmt.Print("Image Version : ")
			dtReader := bufio.NewReader(os.Stdin)
			imgVER, err := dtReader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			imgVER = strings.TrimSpace(imgVER)
			slice := strings.Split(imgVER, " ")
			imgVER = "kubectl delete deployment " + slice[0] + "-" + slice[1] + "-deploy"

			cmd_delete := exec.Command("sh", "-c", imgVER)
			cmd_delete.Stdout = os.Stdout
			if err := cmd_delete.Run(); err != nil {
				panic(err)
			}
			wg.Done()
		}
	default:
		fmt.Println("Check your Command")
		deletePods()
	}
	wg.Wait()
}