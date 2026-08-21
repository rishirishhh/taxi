package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}

	return s.repo.CreateTrip(ctx, t)
}

func (s *Service) GetRoute(
	ctx context.Context,
	pickup, destination *types.Coordinate,
) (*types.OsrmApiResponse, error) {

	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude,
		pickup.Latitude,
		destination.Longitude,
		destination.Latitude,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSRM request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OSRM: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"OSRM returned status code: %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OSRM response: %w", err)
	}

	var routeResp types.OsrmApiResponse

	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse OSRM response: %w", err)
	}

	return &routeResp, nil
}
