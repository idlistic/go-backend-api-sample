package router

import (
	"net/http"

	"github.com/idlistic/go-backend-api-sample/internal/db"
	"github.com/idlistic/go-backend-api-sample/internal/handler"
	"github.com/idlistic/go-backend-api-sample/internal/repository"
)

func New() (http.Handler, func() error, error) {
	database, err := db.Open()
	if err != nil {
		return nil, nil, err
	}

	timeslotRepo := repository.NewTimeslotRepository(database)
	timeslotHandler := handler.NewTimeslotHandler(timeslotRepo)

	branchRepo := repository.NewBranchRepository(database)
	branchHandler := handler.NewBranchHandler(branchRepo)

	orderRepo := repository.NewOrderRepository(database)
	orderHandler := handler.NewOrderHandler(orderRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// GET /timeslots?branch_id=&date=
	mux.HandleFunc("/timeslots", timeslotHandler.List)
	mux.HandleFunc("/branches", branchHandler.List)
	mux.HandleFunc("/orders", orderHandler.Handle)

	mux.HandleFunc("/orders/", orderHandler.Cancel) // for /orders/{id}/cancel

	cleanup := func() error { return database.Close() }
	return withCORS(mux), cleanup, nil
}
