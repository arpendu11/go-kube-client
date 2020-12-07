package main

import (
	"log"
	"os/exec"
	"github.com/urfave/cli"
)

func deploy(c *cli.Context) {
	grep := exec.Command("grep", "grafana")
    ps := exec.Command("kubectl", "get", "pods", "--all-namespaces")

    // Get ps's stdout and attach it to grep's stdin.
    pipe, _ := ps.StdoutPipe()
    defer pipe.Close()

    grep.Stdin = pipe

    // Run ps first.
    ps.Start()

    // Run and get the output of grep.
    res, _ := grep.Output()

    log.Printf("%s\n", string(res))
}