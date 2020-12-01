package main

import (
	"context"
	"fmt"
	"net"
	"log"
	"os"
	"os/signal"

	"github.com/arpendu11/go-kube-client/kubepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/urfave/cli"
)

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

	// execute cli based portion
	app := cli.NewApp()

	app.Name = "kube-client"
	app.Usage = "Using client-go effectively with Kubernetes api"
	app.Version = "1.0"

	app.Commands = []cli.Command{
		{Name: "crud", Usage: "Run CRUD example", Action: crudOperation},
		{Name: "lister", Usage: "Run lister example", Action: lister},
		{Name: "informer", Usage: "Run informer example", Action: informer},
		{Name: "workqueue", Usage: "Run workqueue example", Action: workqueue_example},
		{Name: "deploy", Usage: "Execute Helm Chart to deploy package", Action: deploy},
	}

	app.Run(os.Args)
}