package models

import (
	"database/sql"
	"time"
)

// Reservation represents an instance launch reservation. They are associated with a background
// job system with a particular job with its own ID. The function handlers update the reservation
// Status, Success and FinishedAt attributes until the job is considered finished.
type Reservation struct {
	// Required auto-generated PK.
	ID int64 `db:"id" json:"id"`

	// Provider type. Required.
	Provider ProviderType `db:"provider" json:"provider"`

	// Account ID. Required.
	AccountID int64 `db:"account_id" json:"-"`

	// Time when reservation was made.
	CreatedAt time.Time `db:"created_at" json:"created_at"`

	// Total number of job steps for this reservation.
	Steps int32 `db:"steps" json:"steps"`

	// User-facing step descriptions for each step. Length of StepTitles must be equal to Steps.
	StepTitles []string `db:"step_titles" json:"step_titles"`

	// Active job step for this reservation. See Status for more details.
	Step int32 `db:"step" json:"step"`

	// Textual status of the reservation or error when there was a failure. This is user-facing message,
	// it should not contain error details. See Error field for more details.
	Status string `db:"status" json:"status"`

	// Error contains error message from the last step. When a job returns an error, execution of
	// reservation should stop, Success flag should be set to false and Error will contain the error
	// message.
	Error string `db:"error" json:"error,omitempty"`

	// Time when reservation was finished or nil when it's still processing.
	FinishedAt sql.NullTime `db:"finished_at" json:"finished_at"`

	// Flag indicating success, error or unknown state (NULL). See Status for the actual error.
	Success sql.NullBool `db:"success" json:"success"`
}

type NoopReservation struct {
	Reservation
}

type AWSDetail struct {
	Region string `json:"region"`

	// Optional instance name
	Name *string `json:"name"`

	// AWS Instance type.
	InstanceType string `json:"instance_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int32 `json:"amount"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`

	// PubkeyName on AWS in given region. Found by the EnsurePubkey job.
	PubkeyName string `json:"pubkey_name"`
}

type AWSReservation struct {
	Reservation

	// Pubkey ID.
	PubkeyID int64 `db:"pubkey_id" json:"pubkey_id"`

	// Source ID.
	SourceID string `db:"source_id" json:"source_id"`

	// The ID of the aws reservation which was created or nil when reservation was not yet created.
	AWSReservationID *string `db:"aws_reservation_id" json:"aws_reservation_id"`

	// The ID of the image from which the instance is created. AMI's must be prefixed with 'ami-'.
	ImageID string `db:"image_id" json:"image_id"`

	// Detail information is stored as JSON in DB
	Detail *AWSDetail `db:"detail" json:"detail"`
}

type GCPDetail struct {
	Zone string `json:"zone"`

	// Optional instance name
	Name *string `json:"name"`

	// GCP Machine type.
	MachineType string `json:"machine_type"`

	// Amount of instances to provision of type: Instance type.
	Amount int64 `json:"amount"`

	// Immediately power off the system after initialization
	PowerOff bool `json:"poweroff"`
}

type GCPReservation struct {
	Reservation

	// Pubkey ID.
	PubkeyID int64 `db:"pubkey_id" json:"pubkey_id"`

	// Source ID.
	SourceID string `db:"source_id" json:"source_id"`

	// The ID of the gcp reservation which was created.
	GCPOperationName string `db:"gcp_operation_name"`

	// The ID of the image from which the instance is created.
	ImageID string `json:"image_id"`

	// Detail information is stored as JSON in DB
	Detail *GCPDetail `db:"detail" json:"detail"`
}

type ReservationInstance struct {
	// Reservation ID.
	ReservationID int64 `db:"reservation_id" json:"reservation_id"`

	// Instance ID which has been created on a cloud provider.
	InstanceID string `db:"instance_id" json:"instance_id"`
}
