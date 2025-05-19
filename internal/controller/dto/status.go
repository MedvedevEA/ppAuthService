package dto

import "github.com/google/uuid"

type AddStatus struct {
	Name string `json:"name"`
}
type GetStatus struct {
	StatusId *uuid.UUID `uri:"statusId" validate:"required"`
}
type GetStatuses struct {
}
type UpdateStatus struct {
	StatusId *uuid.UUID `uri:"statusId" validate:"required"`
	Name     *string    `json:"name"`
}
type RemoveStatus struct {
	StatusId *uuid.UUID `uri:"statusId" validate:"required"`
}
