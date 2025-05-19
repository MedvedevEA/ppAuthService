package dto

type GetUsers struct {
	Offset int     `query:"offset" validate:"gte=0"`
	Limit  int     `query:"limit" validate:"gte=0"`
	Name   *string `query:"name" validate:"omitempty"`
}
