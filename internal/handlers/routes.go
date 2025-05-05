package handlers

import (
	"github.com/gorilla/mux"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
)

func NewRouter(h adapters.Handler) *mux.Router {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/ping", h.Ping).Methods("GET")

	rtr.HandleFunc("/place-bet", h.PlaceBet).Methods("POST")

	rtr.HandleFunc("/settle-bets", h.SettleBets).Methods("POST")

	rtr.HandleFunc("/get-account", h.GetAccount).Methods("GET")

	return rtr
}
