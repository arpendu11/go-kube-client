package main

import (
	"context"
	"github.com/arpendu11/go-kube-client/kubepb"
	"fmt"
	"log"
	"io"

	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Hello I'm a client!!")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	c := kubepb.NewDeployServiceClient(conn)

	defer conn.Close()

	doServerStreaming(c)
}

func doServerStreaming(c kubepb.DeployServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	req := &kubepb.DeployRequest{
		DeployManifest: &kubepb.DeployManifest{
			CustomerName: "Cognizant",
			CustomerType: "IT",
			DeploymentType: "Helm",
			Products: []string{"fusion", "recon", "prometheus", "grafana"},
		},
	}
	
	resStream, err := c.Deploy(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling Server Streaming Deploy RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from Deploy: %v", msg.GetResult())
	}
}