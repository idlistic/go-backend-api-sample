package handler

import (
	"net/http"

	"github.com/idlistic/go-backend-api-sample/internal/repository"
)

type BranchHandler struct {
	repo *repository.BranchRepository
}

func NewBranchHandler(repo *repository.BranchRepository) *BranchHandler {
	return &BranchHandler{repo: repo}
}

func (h *BranchHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to query branches",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"count": len(items),
	})
}
