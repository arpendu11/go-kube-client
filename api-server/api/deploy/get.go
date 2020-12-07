package deploy

import (
	"encoding/json"
	"net/http"
	"log"
)

func doGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Request to get deployment!!\n")
	je := json.NewEncoder(w)
	je.Encode(db)
}
