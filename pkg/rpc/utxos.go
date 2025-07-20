package rpc

import (
	"net/http"

	"dogecoin.org/balance-master/pkg/store"
)

type UtxoRoutes struct {
	store *store.Store
}

func HandleUtxoRoutes(store *store.Store, mux *http.ServeMux) {
	hr := &UtxoRoutes{store: store}
	mux.HandleFunc("/utxos", hr.handleUtxos)
}

func (hr *UtxoRoutes) handleUtxos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hr.getUtxos(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (hr *UtxoRoutes) getUtxos(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	utxos, err := hr.store.GetUtxos(address)
	if err != nil {
		http.Error(w, "Error getting utxos", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, utxos)
}
