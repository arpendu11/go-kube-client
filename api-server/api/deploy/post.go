package deploy

import (
	"errors"
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
	// var deployed = false
	if len(db) < 1 {
		deployed, err := installProducts(aDeploy.Products)
		if err != nil {
			log.Fatalf("Failed to install: %v\n", err)
		}
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
						deployed, err1 := installProducts([]string{p})
						if err1 != nil {
							log.Fatalf("Failed to install: %v\n", err1)
						}
						if deployed {
							d.Products = append([]string{p}, d.Products...)
							aDeploy.Status = d.Status
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
			deployed, err2 := installProducts(aDeploy.Products)
			if err2 != nil {
				log.Fatalf("Failed to install: %v\n", err2)
			}
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

func installProducts(products []string) (bool, error) {
	var cmd string
	var installed = false
	for _, p := range products {
		if p == "fusion" {
			cmd = "helm install fusion fusion-helm-chart/. -n $(kubectl get namespaces | grep arcsight | cut -d ' ' -f1)"
			installCheck, err1 := checkInstalled(p)
			if err1 != nil {
				log.Fatalf("Failed to execute commands: %v\n", err1)
				return false, err1
			}
			if !installCheck {
				directDeploy, err := helmInstall(cmd)
				if err != nil {
					log.Fatalf("Failed to execute commands: %v %v\n", cmd, err)
					return false, err
				}
				installed = directDeploy
			} else {
				return true, nil
			}
		} else if p == "recon" {
			cmd = "helm install recon recon-helm-chart/. -n $(kubectl get namespaces | grep arcsight | cut -d ' ' -f1)"
			installCheck, err1 := checkInstalled(p)
			if err1 != nil {
				log.Fatalf("Failed to execute commands: %v\n", err1)
				return false, err1
			}
			if !installCheck {
				directDeploy, err := helmInstall(cmd)
				if err != nil {
					log.Fatalf("Failed to execute commands: %v %v\n", cmd, err)
					return false, err
				}
				installed = directDeploy
			} else {
				return true, nil
			}
		} else if p == "prometheus" {
			cmd = "helm install prometheus prometheus-community/kube-prometheus-stack -n $(kubectl get namespaces | grep arcsight | cut -d ' ' -f1)"
			installCheck, err1 := checkInstalled(p)
			if err1 != nil {
				log.Fatalf("Failed to execute commands: %v\n", err1)
				return false, err1
			}
			if !installCheck {
				directDeploy, err := helmInstall(cmd)
				if err != nil {
					log.Fatalf("Failed to execute commands: %v %v\n", cmd, err)
					return false, err
				}
				installed = directDeploy
			} else {
				return true, nil
			}
		} else if p == "monitoring_dashboard" {
			cmd = "helm install prometheus prometheus-community/kube-prometheus-stack -n $(kubectl get namespaces | grep arcsight | cut -d ' ' -f1)"
			installCheck, err1 := checkInstalled("grafana")
			if err1 != nil {
				log.Fatalf("Failed to execute commands: %v\n", err1)
				return false, err1
			}
			if !installCheck {
				directDeploy, err := helmInstall(cmd)
				if err != nil {
					log.Fatalf("Failed to execute commands: %v %v\n", cmd, err)
					return false, err
				}
				installed = directDeploy
			} else {
				return true, nil
			}
		} else {
			log.Fatalf("Failed to execute commands: Product %s not recognized!\n", p)
			return false, errors.New("Product not recognized")
		}
	}
	return installed, nil
}

func helmInstall(cmd string) (bool, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalf("Failed to execute commands: %v %v\n", cmd, err)
		return false, err
	}
	log.Printf("Deployment Result: {%s}\n", out)
	return true, nil		
}

func checkInstalled(product string) (bool, error) {
	grep := exec.Command("grep", product)
    ps := exec.Command("kubectl", "get", "pods", "--all-namespaces")

    // Get ps's stdout and attach it to grep's stdin.
    pipe, _ := ps.StdoutPipe()
    defer pipe.Close()
    grep.Stdin = pipe

    // Run ps first.
    ps.Start()

    // Run and get the output of grep.
	res, _ := grep.Output()
	
	log.Printf("Result: {%s}\n", string(res))
	if len(string(res)) != 0 {
		return true, nil
	}
	return false, nil
}
