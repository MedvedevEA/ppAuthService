package dto

import "github.com/google/uuid"

type UpdateStatus struct {
	StatusId *uuid.UUID `json:"status_id"`
	Name     *string    `json:"name"`
}
