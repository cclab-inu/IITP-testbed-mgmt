package cluster

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	bar "github.com/cclab.inu/testbed-mgmt/src/bar"
	config "github.com/cclab.inu/testbed-mgmt/src/config"
	"github.com/cclab.inu/testbed-mgmt/src/types"
)

var parsed bool = false
var kubeconfig *string

func isInCluster() bool {
	if _, ok := os.LookupEnv("KUBERNETES_PORT"); ok {
		return true
	}

	return false
}

func ConnectK8sClient() *kubernetes.Clientset {
	if isInCluster() {
		return ConnectInClusterAPIClient()
	}

	return ConnectLocalAPIClient()
}

func ConnectLocalAPIClient() *kubernetes.Clientset {
	if !parsed {
		homeDir := ""
		if h := os.Getenv("HOME"); h != "" {
			homeDir = h
		} else {
			homeDir = os.Getenv("USERPROFILE") // windows
		}

		envKubeConfig := os.Getenv("KUBECONFIG")
		if envKubeConfig != "" {
			kubeconfig = &envKubeConfig
		} else {
			if home := homeDir; home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()
		}

		parsed = true
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	return clientset
}

func ConnectInClusterAPIClient() *kubernetes.Clientset {
	cfg := config.GetCfgCluster()

	host := ""
	port := ""
	token := ""

	if val, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST"); ok {
		host = val
	} else {
		host = cfg.Master
	}

	if val, ok := os.LookupEnv("KUBERNETES_PORT_443_TCP_PORT"); ok {
		port = val
	} else {
		port = "6443"
	}

	read, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	token = string(read)

	// create the configuration by token
	kubeConfig := &rest.Config{
		Host:        "https://" + host + ":" + port,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	if client, err := kubernetes.NewForConfig(kubeConfig); err != nil {
		log.Error().Msg(err.Error())
		return nil
	} else {
		return client
	}
}

func LoadConfigCluster() {
	panic("unimplemented")
}

// =============== //
// == Namespace == //
// =============== //

func GetNamespacesFromK8sClient() []string {
	results := []string{}

	client := ConnectK8sClient()
	if client == nil {
		return results
	}

	// get namespaces from k8s api client
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msg(err.Error())
		return results
	}

	for _, namespace := range namespaces.Items {
		if namespace.Status.Phase != "Active" {
			continue
		}

		results = append(results, namespace.Name)
	}

	return results
}

// ========= //
// == Pod == //
// ========= //

func GetPodsFromK8sClient() []types.Pod {
	results := []types.Pod{}

	client := ConnectK8sClient()
	if client == nil {
		return nil
	}

	// get pods from k8s api client
	pods, err := client.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msg(err.Error())
		return results
	}

	for _, pod := range pods.Items {
		containers := []string{}
		for _, c := range pod.Spec.Containers {
			containers = append(containers, c.Name)
		}

		results = append(results, types.Pod{
			Namespace:  pod.Namespace,
			PodName:    pod.Name,
			PodIP:      pod.Status.PodIP,
			Containers: containers,
		})
	}

	return results
}

// ============= //
// == Cluster == //
// ============= //

func CreateCluster() {
	log.Info().Msg("Creating Cluster...")
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan interface{})
	go bar.StartBar(ch, &wg)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	file := dir + "/../scripts/create_cluster.sh"
	cmd := exec.Command("/bin/bash", file)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	file = config.GetCurrentCfg().Home + "/k8s_init.log"
	b, _ := os.Open(file)
	content, err := ioutil.ReadAll(b)
	if err != nil {
		panic(err)
	}
	output := string(content)

	idxFind := strings.Index(output, "kubeadm join")
	left := strings.LastIndex(output[:idxFind], "\n")
	right := strings.Index(output[idxFind:], "\n")
	line1 := strings.Trim(output[left:idxFind+right], "\\")
	line1 = strings.TrimSpace(line1)

	idxFind = strings.Index(output, "--discovery-token-ca-cert-hash")
	left = strings.LastIndex(output[:idxFind], "\n")
	right = strings.Index(output[idxFind:], "\n")
	line2 := output[left : idxFind+right]
	line2 = strings.TrimSpace(line2)

	join := line1 + " " + line2
	join = strings.Replace(join, "\n", " ", -1)
	join = strings.TrimSpace(join)

	cfgCluster := config.GetCurrentCfg().ConfigCluster
	workers := []types.ConfigWorker{cfgCluster.Worker1, cfgCluster.Worker2}
	for _, worker := range workers {
		sshCient, err := ConnectSSH(worker.IP+":22", worker.SSHID, worker.SSHPW)
		if err != nil {
			log.Err(err)
			os.Exit(1)
		}

		_, err = sshCient.SendCommands("sudo " + join)
		if err != nil {
			log.Err(err)
			os.Exit(1)
		}
	}

	close(ch)
	wg.Wait()
}

func DeleteCluster() {
	log.Info().Msg("Deleting Cluster...")
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan interface{})
	go bar.StartBar(ch, &wg)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	file := dir + "/../scripts/delete_cluster.sh"
	cmd := exec.Command("/bin/bash", file)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	cfgCluster := config.GetCurrentCfg().ConfigCluster
	workers := []types.ConfigWorker{cfgCluster.Worker1, cfgCluster.Worker2}
	for _, worker := range workers {
		sshCient, err := ConnectSSH(worker.IP+":22", worker.SSHID, worker.SSHPW)
		if err != nil {
			panic(err)
		}

		_, err = sshCient.SendCommands("sudo kubeadm reset -f ")
		if err != nil {
			panic(err)
		}
	}

	close(ch)
	wg.Wait()
}
