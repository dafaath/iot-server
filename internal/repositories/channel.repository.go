package repositories

import (
	"context"
	"time"

	"github.com/dafaath/iot-server/internal/entities"
	"github.com/dafaath/iot-server/internal/helper"
)

type ChannelRepository struct{}

func NewChannelRepository() (ChannelRepository, error) {
	return ChannelRepository{}, nil
}

func (c *ChannelRepository) Create(ctx context.Context, tx helper.Querier, payload *entities.ChannelCreate) (entities.Channel, error) {
	channel := entities.Channel{
		Time:          time.Now().UTC(),
		ChannelCreate: *payload,
	}
	sqlStatement := `
	INSERT INTO "channel" (
		time, 
		value, 
		id_sensor)
	VALUES ($1, $2, $3)`
	_, err := tx.Exec(ctx, sqlStatement, channel.Time, channel.Value, channel.IdSensor)
	if err != nil {
		return channel, err
	}

	return channel, nil
}
