package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/cilium/cilium/api/v1/flow"
	"github.com/cilium/cilium/api/v1/observer"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"

	cluster "github.com/cclab.inu/testbed-mgmt/src/cluster"
	"github.com/cclab.inu/testbed-mgmt/src/config"
	"github.com/cclab.inu/testbed-mgmt/src/types"
	cilium "github.com/cilium/cilium/api/v1/flow"
	pb "github.com/kubearmor/KubeArmor/protobuf"
)

// ======================= //
// == Global Variables  == //
// ======================= //

/*
var HubbleURL string
var HubblePort string

var KubeArmorURL string
var KubeArmorPort string
*/

var SystemStopChan chan struct{}
var NetworkStopChan chan struct{}
var AppStopChan chan struct{}

// init Function
func init() {
	SystemStopChan = make(chan struct{})
	NetworkStopChan = make(chan struct{})
	AppStopChan = make(chan struct{})
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

// ========================= //
// == Cilium Hubble Relay == //
// ========================= //

func ConnectHubbleRelay(cfg types.ConfigCiliumHubble) *grpc.ClientConn {
	addr := net.JoinHostPort(cfg.HubbleURL, cfg.HubblePort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	log.Info().Msg("dialed for hubble relay:" + addr)
	return conn
}

var HubbleRelayStarted = false

func StartHubbleRelay(StopChan chan struct{}, cfg types.ConfigCiliumHubble) {
	if HubbleRelayStarted {
		return
	}

	conn := ConnectHubbleRelay(cfg)
	if conn == nil {
		log.Error().Msg("ConnectHubbleRelay() failed")
		return
	}
	HubbleRelayStarted = true

	defer func() {
		log.Info().Msg("hubble relay stream rcvr returning")
		HubbleRelayStarted = false
		_ = conn.Close()
	}()

	client := observer.NewObserverClient(conn)

	req := &observer.GetFlowsRequest{
		Follow: true,
		Whitelist: []*cilium.FlowFilter{
			{
				TcpFlags: []*flow.TCPFlags{
					{SYN: true},
				},
			},
			{
				Protocol: []string{"udp"},
				Reply:    []bool{false},
			},
			{
				Protocol: []string{"icmp", "http", "dns"},
			},
		},
		Blacklist: []*cilium.FlowFilter{
			{
				TcpFlags: []*flow.TCPFlags{
					{ACK: true},
				},
			},
		},
	}

	stream, err := client.GetFlows(context.Background(), req)
	if err != nil {
		log.Error().Msg("Unable to stream network flow: " + err.Error())
		return
	}

	for {
		select {
		case <-StopChan:
			return

		default:
			res, err := stream.Recv()
			if err != nil {
				log.Error().Msg("Cilium network flow stream stopped: " + err.Error())
				return
			}

			switch r := res.ResponseTypes.(type) {
			case *observer.GetFlowsResponse_Flow:
				flow := r.Flow
				b, _ := flow.MarshalJSON()

				fmt.Println(jsonPrettyPrint(string(b))) // get cilium logs
			}
		}
	}
}

// ===================== //
// == KubeArmor Relay == //
// ===================== //

func ConnectKubeArmorRelay(cfg types.ConfigKubeArmorRelay) *grpc.ClientConn {
	addr := net.JoinHostPort(cfg.KubeArmorURL, cfg.KubeArmorPort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("err connecting kubearmor relay. " + err.Error())
		return nil
	}

	log.Info().Msg("connected to kubearmor relay " + addr)
	return conn
}

var KubeArmorRelayStarted = false

func StartKubeArmorRelay(StopChan chan struct{}, cfg types.ConfigKubeArmorRelay) {
	if KubeArmorRelayStarted {
		log.Info().Msg("kubearmor relay already started")
		return
	}
	KubeArmorRelayStarted = true
	conn := ConnectKubeArmorRelay(cfg)

	client := pb.NewLogServiceClient(conn)
	req := pb.RequestMessage{}
	req.Filter = "all"

	defer func() {
		log.Info().Msg("watchlogs returning")
		KubeArmorRelayStarted = false
		_ = conn.Close()
	}()

	stream, err := client.WatchLogs(context.Background(), &req)
	if err != nil {
		log.Error().Msg("unable to stream systems logs: " + err.Error())
		return
	}

	for {
		select {
		case <-StopChan:
			return

		default:
			res, err := stream.Recv()
			if err != nil {
				log.Error().Msg("watch logs stream stopped: " + err.Error())
				return
			}

			log.Info().Msg(res.GetData())
		}
	}
}

// ======================= //
// == Kubernetes System == //
// ======================= //

func StartPodLogs(StopChan chan struct{}, pod types.Pod) {
	var since int64 = 1
	podLogOptions := v1.PodLogOptions{
		Follow:       true,
		SinceSeconds: &since,
		Container:    pod.Containers[0],
	}

	cli := cluster.ConnectK8sClient()
	podLogRequest := cli.CoreV1().Pods(pod.Namespace).GetLogs(pod.PodName, &podLogOptions)

	stream, err := podLogRequest.Stream(context.Background())
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	defer stream.Close()

	for {
		select {
		case <-StopChan:
			return

		default:
			buf := make([]byte, 1000)
			numBytes, err := stream.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Error().Msg(err.Error())
				return
			}

			log.Info().Msg(string(buf[:numBytes])) // get k8s logs
		}
	}
}

// ================ //
// == print-logs == //
// ================ //

func PrintLogs() {
	var wg sync.WaitGroup
	wg.Add(1)

	pods := cluster.GetPodsFromK8sClient()

	switch os.Args[2] {
	case "all":
		log.Info().Msg("** Network Level Log **\n")
		go StartHubbleRelay(NetworkStopChan, config.GetCfgCiliumHubble())
		time.Sleep(time.Second * 5)
		log.Info().Msg("\n")

		log.Info().Msg("** System Level Log **\n")
		go StartKubeArmorRelay(SystemStopChan, config.GetCfgKubeArmor())
		time.Sleep(time.Second * 5)
		log.Info().Msg("\n")

		log.Info().Msg("** Application Level Log **\n")
		for _, pd := range pods {
			go StartPodLogs(AppStopChan, pd)
		}

		time.Sleep(time.Second * 5)
		wg.Done()

	case "network":
		log.Info().Msg("** Network Level Log **")
		go StartHubbleRelay(NetworkStopChan, config.GetCfgCiliumHubble())
		time.Sleep(time.Second * 5)
		wg.Done()

	case "system":
		log.Info().Msg("** System Level Log **")
		go StartKubeArmorRelay(SystemStopChan, config.GetCfgKubeArmor())
		time.Sleep(time.Second * 5)
		wg.Done()

	case "app":
		log.Info().Msg("** Application Level Log **")
		for _, pd := range pods {
			go StartPodLogs(AppStopChan, pd)
		}

		time.Sleep(time.Second * 5)
		wg.Done()

	default:
		log.Info().Msg("Check your Command")
	}

}
