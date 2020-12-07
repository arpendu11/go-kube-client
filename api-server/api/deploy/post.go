package deploy

import (
	"log"
	"encoding/json"
	"net/http"
	"reflect"
	"os/exec"
)

func doPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jd := json.NewDecoder(r.Body)

	aDeploy := &DeployManifest{}
	err := jd.Decode(aDeploy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Request to post deployment: %v\n", aDeploy)

	// start of protected code changes
	lock.Lock()
	var deployed = false	
	if len(db) < 1 {
		deployed = installProducts(aDeploy.Products)
		if deployed {
			aDeploy.Status = "Running"
			nextDeployID++
			aDeploy.ID = nextDeployID
			db = append(db, aDeploy)
		} else {
			aDeploy.Status = "Failed"
			log.Fatal("Failed to install the products. Please try again later!")
		}		
	} else {
		var match = false
		for _, d := range db {
			if aDeploy.CustomerName == d.CustomerName {
				for _, p := range aDeploy.Products {
					if !itemExists(d.Products, p) {
						deployed = installProducts([]string{p})
						if deployed {
							d.Products = append([]string{p}, d.Products...)
						} else {
							aDeploy.Status = "Failed"
							log.Fatal("Failed to install the products. Please try again later!")
						}
					}
				}
				match = true
				aDeploy.ID = d.ID
			}
		}
		if !match {
			deployed = installProducts(aDeploy.Products)
			if deployed {
				aDeploy.Status = "Running"
				nextDeployID++
				aDeploy.ID = nextDeployID
				db = append(db, aDeploy)
			} else {
				aDeploy.Status = "Failed"
				log.Fatal("Failed to install the products. Please try again later!")
			}
		}
	}
	// end protected code changes
	lock.Unlock()

	respDeploy := DeployResponse{ID: aDeploy.ID, CustomerName: aDeploy.CustomerName, Status: aDeploy.Status}
	je := json.NewEncoder(w)
	je.Encode(respDeploy)
}

func itemExists(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func installProducts(products []string) bool {
	var cmdList []string
	var grepList []string
	for _, p := range products {
		if p == "fusion" {
			cmdList = []string{"helm", "install", "fusion", "fusion-helm-chart/.", "-n", "arcsight-installer-9qe5i"}
		} else if p == "recon" {
			cmdList = []string{"helm", "install", "recon", "recon-helm-chart/.", "-n", "arcsight-installer-9qe5i"}
		} else if p == "prometheus" {
			cmdList = []string{"kubectl", "get", "pods", "--all-namespaces"}
			grepList = []string{"grep", "prometheus"}
		} else if p == "monitoring_dashboard" {
			cmdList = []string{"kubectl", "get", "pods", "--all-namespaces"}
			grepList = []string{"grep", "grafana"}
		} else {
			log.Fatalf("Failed to execute commands: Product %s not recognized!\n", p)
			return false
		}
	}

	if len(cmdList) == 6 {
		cmd, err := exec.Command(cmdList[0], cmdList[1], cmdList[2], cmdList[3], cmdList[4], cmdList[5]).Output()
		if err != nil {
			log.Fatalf("Failed to execute commands: %v\n", err)
			return false
		}
		log.Printf("Deployment Result: {%s}\n", cmd)
	} else {

		grep := exec.Command(grepList[0], grepList[1])
        kube := exec.Command(cmdList[0], cmdList[1], cmdList[2], cmdList[3])

        // Get ps's stdout and attach it to grep's stdin.
        pipe, _ := kube.StdoutPipe()
        defer pipe.Close()

        grep.Stdin = pipe

        // Run ps first.
        kube.Start()

        // Run and get the output of grep.
        res, _ := grep.Output()

		log.Printf("Result: {%s}\n", string(res))
		
		if len(string(res)) != 0 {
			log.Printf("Deployment already installed: {monitoring}\n")
		} else {
			cmd1, err1 := exec.Command("helm", "install", "prometheus", "prometheus-community/kube-prometheus-stack", "-n", "arcsight-installer-9qe5i").Output()
			if err1 != nil {
				log.Fatalf("Failed to execute commands: %v\n", err1)
				return false
			}
			log.Printf("Deployment Result: {%s}\n", cmd1)
		}
	}
	return true
}
