package deploy

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type DeployResponse struct{
	ID				uint64	`json:"id,omitempty"`
	CustomerName	string	`json:"customerName,omitempty"`
	Status			string 	`json:"status,omitempty"`
}
type DeployManifest struct {
	ID				uint64		`json:"id,omitempty"`
	CustomerName	string		`json:"customerName,omitempty"`
	CustomerType	string		`json:"customerType,omitempty"`
	DeploymentType	string		`json:"deploymentType,omitempty"`
	Products		[]string	`json:"products,omitempty"`
	Status			string		`json:"status,omitempty"`
}

var db = []*DeployManifest{}
var nextDeployID uint64
var lock sync.Mutex

func (u *DeployResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case http.MethodGet:
		doGet(w, r)
	case http.MethodPost:
		doPost(w, r)
	case http.MethodDelete:
		doDelete(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unsupported method '%v' to %v\n", r.Method, r.URL)
		log.Printf("Unsupported method '%v' to %v\n", r.Method, r.URL)
	}
}
