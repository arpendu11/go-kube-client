package main

import (
	// "bytes"
	"fmt"
	"log"
	"os/exec"
	// "strings"
	"github.com/urfave/cli"
)

func deploy(c *cli.Context) {
	cmd, err := exec.Command("helm install prometheus prometheus-community/kube-prometheus-stack -n arcsight-installer-9qe5i").Output()
	// cmd.Stdin = strings.NewReader("some input")
	// var out bytes.Buffer
	// cmd.Stdout = &out
	// err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n",cmd)
}