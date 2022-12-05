package consumer

import (
	"context"
	"io"
	"net"

	"github.com/cilium/cilium/api/v1/flow"
	"github.com/cilium/cilium/api/v1/observer"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"

	cluster "github.com/cclab.inu/testbed-mgmt/src/cluster"
	cilium "github.com/cilium/cilium/api/v1/flow"
	pb "github.com/kubearmor/KubeArmor/protobuf"
)

// ======================= //
// == Global Variables  == //
// ======================= //

var HubbleURL string
var HubblePort string

var KubeArmorURL string
var KubeArmorPort string

func init() {
	HubbleURL = "10.109.38.101"
	HubblePort = "80"

	KubeArmorURL = "10.111.137.209"
	KubeArmorPort = "32767"
}

// ========================= //
// == Cilium Hubble Relay == //
// ========================= //

func ConnectHubbleRelay() *grpc.ClientConn {
	addr := net.JoinHostPort(HubbleURL, HubblePort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error().Err(err)
		return nil
	}

	log.Info().Msg("dialed for hubble relay:" + addr)
	return conn
}

var HubbleRelayStarted = false

func StartHubbleRelay(StopChan chan struct{}) {
	if HubbleRelayStarted {
		return
	}

	conn := ConnectHubbleRelay()
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

				println(string(b)) // get cilium logs
			}
		}
	}
}

// ===================== //
// == KubeArmor Relay == //
// ===================== //

func ConnectKubeArmorRelay() *grpc.ClientConn {
	addr := net.JoinHostPort(KubeArmorURL, KubeArmorPort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Error().Msg("err connecting kubearmor relay. " + err.Error())
		return nil
	}

	log.Info().Msg("connected to kubearmor relay " + addr)
	return conn
}

var KubeArmorRelayStarted = false

func StartKubeArmorRelay(StopChan chan struct{}) {
	if KubeArmorRelayStarted {
		log.Info().Msg("kubearmor relay already started")
		return
	}
	KubeArmorRelayStarted = true
	conn := ConnectKubeArmorRelay()

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

			println(res.GetData()) // get kubearmor log
		}
	}
}

// ======================= //
// == Kubernetes System == //
// ======================= //

func StartPodLogs(StopChan chan struct{}, pod, namespace string) {
	var since int64 = 1
	podLogOptions := v1.PodLogOptions{
		Follow:       true,
		SinceSeconds: &since,
	}

	cli := cluster.ConnectK8sClient()
	podLogRequest := cli.CoreV1().Pods(namespace).GetLogs(pod, &podLogOptions)

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

			println(string(buf[:numBytes])) // get k8s logs
		}
	}
}