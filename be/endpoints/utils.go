package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//AddRouterEndpoints add the actual endpoints for api
func AddRouterEndpoints(r *mux.Router) *mux.Router {
	r.HandleFunc("/cal-engagement-rate/{username}", calEngagementRate).Methods("GET")
	return r
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}
