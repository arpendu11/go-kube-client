package deploy

import (
	"log"
	"encoding/json"
	"net/http"
	"reflect"
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
	aDeploy.Status = "Running"	
	if len(db) < 1 {
		nextDeployID++
		aDeploy.ID = nextDeployID
		db = append(db, aDeploy)
	} else {
		var match = false
		for _, d := range db {
			if aDeploy.CustomerName == d.CustomerName {
				for _, p := range aDeploy.Products {
					if !itemExists(d.Products, p) {
						d.Products = append([]string{p}, d.Products...)
					}
				}
				match = true
				aDeploy.ID = d.ID
			}
		}
		if !match {
			nextDeployID++
			aDeploy.ID = nextDeployID
			db = append(db, aDeploy)
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