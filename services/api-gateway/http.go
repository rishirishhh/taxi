package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
)

type tripHandler struct {
	tripService *grpc_clients.TripServiceClient
}

func NewTripHandler(tripService *grpc_clients.TripServiceClient) *tripHandler {
	return &tripHandler{tripService: tripService}
}

func (h *tripHandler) handleTripStart(w http.ResponseWriter, r *http.Request) {
	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	trip, err := h.tripService.Client.CreateTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to start a trip: %v", err)
		http.Error(w, "Failed to start trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: trip}
	writeJSON(w, http.StatusCreated, response)
}

func (h *tripHandler) handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tripPreview, err := h.tripService.Client.PreviewTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("Failed to preview a trip: %v", err)
		http.Error(w, "Failed to preview trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreview}
	writeJSON(w, http.StatusCreated, response)
}
