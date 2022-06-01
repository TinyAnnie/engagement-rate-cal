package endpoints

import (
	"github.com/TinyAnnie/engagement-rate-cal/be/services"
	"github.com/gorilla/mux"
	"net/http"
)

//calEngagementRate will return all the posts actually in the array
func calEngagementRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars["username"]
	if !ok {
		http.Error(w, "Invalid param", http.StatusBadRequest)
		return
	}
	rate := services.CalEngagementRate2(username)
	sendJSONResponse(w, rate)
}