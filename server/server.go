package main

import (
	"time"
	"fmt"
	"net"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/arpendu11/go-kube-client/kubepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {}

func (*server) Deploy(req *kubepb.DeployRequest, stream kubepb.DeployService_DeployServer) error {
	fmt.Printf("Deploy all call was invoked with %v\n", req)
	firstName := req.GetDeployManifest().GetCustomerName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &kubepb.DeployResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("Welcome to the world of kube-client - A CLI based client app to call Kubernetes APIs !!")
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Deploy Service started!!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	kubepb.RegisterDeployServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Ctrl+c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("End of Program")
}