package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/idlistic/go-backend-api-sample/internal/repository"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

type CreateOrderRequest struct {
	BranchID     int64  `json:"branch_id"`
	TimeslotID   int64  `json:"timeslot_id"`
	CustomerName string `json:"customer_name"`
}

func (h *OrderHandler) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.List(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "method not allowed",
		})
	}
}
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid json body",
		})
		return
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)

	if req.BranchID <= 0 || req.TimeslotID <= 0 || req.CustomerName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "branch_id, timeslot_id, customer_name are required",
		})
		return
	}

	order, err := h.repo.CreateWithTimeslotReservation(r.Context(), req.BranchID, req.TimeslotID, req.CustomerName)
	if err != nil {
		switch err {
		case repository.ErrTimeslotNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "timeslot not found"})
			return
		case repository.ErrTimeslotInactive:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "timeslot inactive"})
			return
		case repository.ErrTimeslotFullyBooked:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "timeslot fully booked"})
			return
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed to create order"})
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"order": order,
	})
}

func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	// Expect: PATCH /orders/{id}/cancel
	if r.Method != http.MethodPatch {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "method not allowed",
		})
		return
	}

	path := r.URL.Path // e.g. /orders/123/cancel
	const prefix = "/orders/"
	const suffix = "/cancel"

	if len(path) <= len(prefix)+len(suffix) || path[:len(prefix)] != prefix || path[len(path)-len(suffix):] != suffix {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	idStr := path[len(prefix) : len(path)-len(suffix)]
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || orderID <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid order id",
		})
		return
	}

	order, err := h.repo.CancelAndReleaseTimeslot(r.Context(), orderID)
	if err != nil {
		switch err {
		case repository.ErrOrderNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "order not found"})
			return
		case repository.ErrOrderNotCancellable:
			writeJSON(w, http.StatusConflict, map[string]any{"error": "order not cancellable"})
			return
		default:
			// if timeslot missing etc.
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "failed to cancel order"})
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"order": order,
	})
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
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

	if _, err := time.Parse("2006-01-02", date); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "date must be YYYY-MM-DD",
		})
		return
	}

	items, err := h.repo.ListByBranchAndDate(r.Context(), branchID, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to query orders",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"branch_id": branchID,
		"date":      date,
		"count":     len(items),
		"items":     items,
	})
}
