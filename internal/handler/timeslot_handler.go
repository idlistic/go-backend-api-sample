package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/idlistic/go-backend-api-sample/internal/repository"
)

type TimeslotHandler struct {
	repo *repository.TimeslotRepository
}

func NewTimeslotHandler(repo *repository.TimeslotRepository) *TimeslotHandler {
	return &TimeslotHandler{repo: repo}
}

func (h *TimeslotHandler) List(w http.ResponseWriter, r *http.Request) {
	branchIDStr := r.URL.Query().Get("branch_id")
	date := r.URL.Query().Get("date")

	if branchIDStr == "" || date == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "branch_id and date are required",
		})
		return
	}

	branchID, err := strconv.ParseInt(branchIDStr, 10, 64)
	if err != nil || branchID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "branch_id must be a positive integer",
		})
		return
	}

	// validate date format YYYY-MM-DD
	if _, err := time.Parse("2006-01-02", date); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "date must be YYYY-MM-DD",
		})
		return
	}

	items, err := h.repo.ListByBranchAndDate(r.Context(), branchID, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to query timeslots",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"count": len(items),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
