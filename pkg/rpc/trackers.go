package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/dogecoinfoundation/balance-master/pkg/store"
)

type TrackerRoutes struct {
	store *store.Store
}

func HandleTrackerRoutes(store *store.Store, mux *http.ServeMux) {
	hr := &TrackerRoutes{store: store}

	mux.HandleFunc("/trackers", hr.handleTrackers)
}

func (hr *TrackerRoutes) handleTrackers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		hr.postTrackers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (hr *TrackerRoutes) postTrackers(w http.ResponseWriter, r *http.Request) {
	var request PostTrackersRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tracker := store.Tracker{
		Address: request.Address,
	}

	err = hr.store.SaveTracker(tracker)
	if err != nil {
		http.Error(w, "Error saving tracker", http.StatusInternalServerError)
		return
	}

	response := PostTrackersResponse{
		ID: tracker.ID,
	}

	respondJSON(w, http.StatusOK, response)
}
