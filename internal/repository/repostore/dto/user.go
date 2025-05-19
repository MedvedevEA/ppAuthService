package dto

import "github.com/google/uuid"

type AddUser struct {
	UserId *uuid.UUID `json:"user_id"`
	Name   string     `json:"name"`
}
type GetUsers struct {
	Offset int     `json:"offset"`
	Limit  int     `json:"limit"`
	Name   *string `json:"name"`
}
