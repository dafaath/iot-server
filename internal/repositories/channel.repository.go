package repositories

import (
	"context"
	"strconv"
	"time"

	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
)

type ChannelRepository struct{}

func NewChannelRepository() (ChannelRepository, error) {
	return ChannelRepository{}, nil
}

func (c *ChannelRepository) Create(ctx context.Context, tx helper.Querier, payload *entities.ChannelCreate) error {
	time := time.Now().UTC()
	sqlStatement := `
	INSERT INTO "channel" (
		time, 
		value, 
		id_node)
	VALUES ($1, $2, $3)`
	_, err := tx.Exec(ctx, sqlStatement, time, payload.Value, payload.IdNode)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChannelRepository) GetNodeChannel(ctx context.Context, tx helper.Querier, idNode int, limit int) ([]entities.Channel, error) {
	channels := []entities.Channel{}
	sqlStatement := `SELECT time, value, id_node FROM "channel" WHERE id_node=$1`
	if limit >= 0 {
		sqlStatement += " LIMIT " + strconv.Itoa(limit)
	}
	rows, err := tx.Query(ctx, sqlStatement, idNode)
	if err != nil {
		return channels, err
	}
	defer rows.Close()

	for rows.Next() {
		var channel entities.Channel
		err := rows.Scan(
			&channel.Time,
			&channel.Value,
			&channel.IdNode,
		)
		if err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}
	if err := rows.Err(); err != nil {
		return channels, err
	}

	return channels, nil

}
