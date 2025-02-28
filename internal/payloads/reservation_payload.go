package payloads

import (
	"net/http"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/go-chi/render"
)

// ReservationRequest is empty, account comes in HTTP header and
// provider type in HTTP URL. All other fields are auto-generated.

type GenericReservationResponsePayload struct {
	ID int64 `json:"id"`

	// Provider type. Required.
	Provider int `json:"provider"`

	// Time when reservation was made.
	CreatedAt time.Time `json:"created_at"`

	// Total number of job steps for this reservation.
	Steps int32 `json:"steps"`

	// User-facing step descriptions for each step. Length of StepTitles must be equal to Steps.
	StepTitles []string `json:"step_titles"`

	// Active job step for this reservation. See Status for more details.
	Step int32 `json:"step"`

	// Textual status of the reservation or error when there was a failure
	Status string `json:"status"`

	// Error message when reservation was not successful. Only set when Success if false.
	Error string `json:"error"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt *time.Time `json:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success *bool `json:"success"`
}

type AWSReservationResponsePayload struct {
	ID int64 `json:"reservation_id"`

	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	// AWS region.
	Region string `json:"region"`

	// AWS Instance type.
	InstanceType string `json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 `json:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id"`

	// The ID of the aws reservation which was created, or missing if not created yet.
	AWSReservationID string `json:"aws_reservation_id,omitempty"`

	// Optional name of the instance(s).
	Name string `json:"name"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`

	// Instance IDs, only present for finished reservations
	Instances []string `json:"instances,omitempty"`
}

type GCPReservationResponsePayload struct {
	ID int64 `json:"reservation_id"`

	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	// GCP zone.
	Zone string `json:"zone"`

	// GCP Machine type.
	MachineType string `json:"machine_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 `json:"amount"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id"`

	// The name of the gcp operation which was created.
	GCPOperationName string `json:"gcp_operation_name,omitempty"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`
}

type NoopReservationResponsePayload struct {
	ID int64 `json:"reservation_id"`
}

type AWSReservationRequestPayload struct {
	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	// AWS region.
	Region string `json:"region"`

	// Optional name of the instance(s).
	Name string `json:"name"`

	// AWS Instance type.
	InstanceType string `json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 ` json:"amount"`

	// Image Builder UUID of the image that should be launched. AMI's must be prefixed with 'ami-'.
	ImageID string `json:"image_id"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`
}

type GCPReservationRequestPayload struct {
	// Pubkey ID.
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID.
	SourceID string `json:"source_id"`

	// GCP zone.
	Zone string `json:"zone"`

	// GCP Machine type.
	MachineType string `json:"machine_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 ` json:"amount"`

	// Image Builder UUID of the image that should be launched.
	ImageID string `json:"image_id"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`
}

func (p *GenericReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *AWSReservationRequestPayload) Bind(_ *http.Request) error {
	return nil
}

func (p *AWSReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *GCPReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (p *NoopReservationResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewReservationResponse(reservation *models.Reservation) render.Renderer {
	return reservationResponseMapper(reservation)
}

func (p *GCPReservationRequestPayload) Bind(_ *http.Request) error {
	return nil
}

func NewAWSReservationResponse(reservation *models.AWSReservation, instances []*models.ReservationInstance) render.Renderer {
	instanceIds := make([]string, len(instances))
	for iter, inst := range instances {
		instanceIds[iter] = inst.InstanceID
	}

	response := AWSReservationResponsePayload{
		PubkeyID:     reservation.PubkeyID,
		ImageID:      reservation.ImageID,
		SourceID:     reservation.SourceID,
		Region:       reservation.Detail.Region,
		Amount:       reservation.Detail.Amount,
		InstanceType: reservation.Detail.InstanceType,
		ID:           reservation.ID,
		Name:         StringNullToEmpty(reservation.Detail.Name),
		PowerOff:     reservation.Detail.PowerOff,
		Instances:    instanceIds,
	}
	if reservation.AWSReservationID != nil {
		response.AWSReservationID = *reservation.AWSReservationID
	}
	return &response
}

func NewGCPReservationResponse(reservation *models.GCPReservation) render.Renderer {
	response := GCPReservationResponsePayload{
		PubkeyID:         reservation.PubkeyID,
		ImageID:          reservation.ImageID,
		SourceID:         reservation.SourceID,
		Zone:             reservation.Detail.Zone,
		Amount:           reservation.Detail.Amount,
		MachineType:      reservation.Detail.MachineType,
		GCPOperationName: reservation.GCPOperationName,
		ID:               reservation.ID,
		PowerOff:         reservation.Detail.PowerOff,
	}
	return &response
}

func NewNoopReservationResponse(reservation *models.NoopReservation) render.Renderer {
	return &NoopReservationResponsePayload{
		ID: reservation.ID,
	}
}

func NewReservationListResponse(reservations []*models.Reservation) []render.Renderer {
	list := make([]render.Renderer, len(reservations))
	for i, reservation := range reservations {
		list[i] = reservationResponseMapper(reservation)
	}
	return list
}

func reservationResponseMapper(reservation *models.Reservation) *GenericReservationResponsePayload {
	var finishedAt *time.Time
	if reservation.FinishedAt.Valid {
		finishedAt = &reservation.FinishedAt.Time
	}
	var success *bool
	if reservation.Success.Valid {
		success = &reservation.Success.Bool
	}
	return &GenericReservationResponsePayload{
		ID:         reservation.ID,
		Provider:   int(reservation.Provider),
		CreatedAt:  reservation.CreatedAt,
		FinishedAt: finishedAt,
		Status:     reservation.Status,
		Success:    success,
		Steps:      reservation.Steps,
		Step:       reservation.Step,
		StepTitles: reservation.StepTitles,
		Error:      reservation.Error,
	}
}
