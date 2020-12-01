package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	fmt.Println("Welcome to the world of kube-client - A CLI based client app to call Kubernetes APIs !!")

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