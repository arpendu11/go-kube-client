package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"net/http"

	"github.com/arpendu11/go-kube-client/api-server/api/deploy"
)

func main() {
	fmt.Println("Welcome to the world of kube-client - A Server Sent Events based server to deploy products/apps in Kubernetes !!")
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Deploy Service started!!")
	// register RESTful endpoint handler for '/deploy/
	http.Handle("/deploy/", &deploy.DeployResponse{})
	server := &http.Server{Addr: ":50052", Handler: nil}

    go func() {
        if err := server.ListenAndServe(); err != nil {
            log.Fatalf("Failed to serve: %v\n", err)
        }
    }()

    // Wait for Ctrl+c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	if err1 := server.Shutdown(context.TODO()); err1 != nil {
        log.Fatalf("Failed to shutdown the server: %v\n", err1)
    }
	fmt.Println("End of Program")
}