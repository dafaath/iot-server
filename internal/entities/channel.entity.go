package entities

import "time"

type Channel struct {
	Time   time.Time         `json:"time" validate:"required"`
	Value  []JsonNullFloat64 `json:"value" validate:"required"`
	IdNode int               `json:"id_node" validate:"required"`
}

type ChannelCreate struct {
	Value  string `json:"value" validate:"required"`
	IdNode int    `json:"id_node" validate:"required"`
}
