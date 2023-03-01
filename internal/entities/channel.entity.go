package entities

import "time"

type Channel struct {
	Time time.Time `json:"time" validate:"required"`
	ChannelCreate
}

type ChannelCreate struct {
	Value    float64 `json:"value" validate:"required"`
	IdSensor int     `json:"id_sensor" validate:"required"`
}
