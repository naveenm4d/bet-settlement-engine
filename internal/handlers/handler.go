package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/naveenm4d/bet-settlement-engine/internal/core/adapters"
	"github.com/naveenm4d/bet-settlement-engine/pkg/constants"
	"github.com/naveenm4d/bet-settlement-engine/pkg/entities"
)

type server struct {
	betsValidator    adapters.BetsValidator
	resultsValidator adapters.ResultsValidator
	betsService      adapters.BetService
	resultsService   adapters.ResultsService
	accountsService  adapters.AccountsService
	logger           *zap.SugaredLogger
}

func NewHandler(
	betsValidator adapters.BetsValidator,
	resultsValidator adapters.ResultsValidator,
	betsService adapters.BetService,
	resultsService adapters.ResultsService,
	accountsService adapters.AccountsService,
	logger *zap.SugaredLogger,
) adapters.Handler {
	srv := &server{
		betsValidator:    betsValidator,
		resultsValidator: resultsValidator,
		betsService:      betsService,
		resultsService:   resultsService,
		accountsService:  accountsService,
		logger:           logger,
	}

	return srv
}

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *server) PlaceBet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("PlaceBet: bet read body failed | Err - %s", err)
	}

	betInput := &entities.Bet{}

	if err = json.Unmarshal(bodyBytes, betInput); err != nil {
		s.logger.Error("PlaceBet: bet unmarshal failed | Err - %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	if err = s.betsValidator.ValidateBet(ctx, betInput); err != nil {
		s.logger.Errorf("PlaceBet: validation failed! | UserID - %s | Err - %s", betInput.UserID, err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	bet, errB := s.betsService.PlaceBet(ctx, betInput)
	if errB != nil {
		s.logger.Error("PlaceBet: bet placement failed! | UserID - %s | Err - %s", betInput.UserID, errB)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: errB.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bet)
}

func (s *server) SettleBets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("SettleBets: bet read body failed | Err - %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	event := &entities.Event{}

	if err = json.Unmarshal(bodyBytes, event); err != nil {
		s.logger.Error("SettleBets: bet unmarshal failed | Err - %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	if err = s.resultsValidator.ValidateEvent(event); err != nil {
		s.logger.Error("SettleBets: event validation failed | EventID - %s | Err - %s", event.EventID, err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	if err = s.resultsService.SettleBetsForEvent(event.EventID, event.ResultStatus); err != nil {
		s.logger.Error("SettleBets: bet settlement failed | EventID - %s | Err - %s", event.EventID, err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entities.JsonResponse{Success: "true"})
}

func (s *server) GetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.Error("GetAccount: read body failed | Err - %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	var data map[string]any

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		s.logger.Error("GetAccount: unmarshal failed | Err - %s", err)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: err.Error()})
		return
	}

	userID, ok := data["user_id"].(string)
	if !ok {
		s.logger.Error("GetAccount: invalid user id")

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: constants.ErrInvalidUserID.Error()})
		return
	}

	account, errA := s.accountsService.GetAccount(userID)
	if errA != nil {
		s.logger.Error("GetAccount: bet settlement failed | Err - %s", errA)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entities.JsonResponse{Success: "false", Message: errA.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}
