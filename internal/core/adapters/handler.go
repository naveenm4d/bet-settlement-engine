package adapters

import "net/http"

type Handler interface {
	Ping(w http.ResponseWriter, r *http.Request)

	PlaceBet(w http.ResponseWriter, r *http.Request)
	SettleBets(w http.ResponseWriter, r *http.Request)
	GetAccount(w http.ResponseWriter, r *http.Request)
}
