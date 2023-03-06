package entities

import "time"

type Channel struct {
	Time   time.Time `json:"time" validate:"required"`
	Value  []float64 `json:"value" validate:"required"`
	IdNode int       `json:"id_node" validate:"required"`
}

type ChannelCreate struct {
	Value  []float64 `json:"value" validate:"required"`
	IdNode int       `json:"id_node" validate:"required"`
}
