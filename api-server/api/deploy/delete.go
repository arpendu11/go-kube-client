package deploy

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

// doDelete deletes a deploy manifest from the db using the path '/deploy/id', eg: /deploy/2
func doDelete(w http.ResponseWriter, r *http.Request) {

	// get the deploy ID from the path
	fields := strings.Split(r.URL.String(), "/")
	id, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Request to delete deployment: %v", id)

	// start of protected code changes
	lock.Lock()
	var tmp = []*DeployManifest{}
	for _, d := range db {
		if id == d.ID {
			continue
		}
		tmp = append(tmp, d)
	}
	db = tmp
	// end protected code changes
	lock.Unlock()
}