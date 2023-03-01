package repositories

import (
	"context"
	"fmt"

	"github.com/dafaath/iot-server/internal/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type SensorRepository struct{}

func NewSensorRepository() (SensorRepository, error) {
	return SensorRepository{}, nil
}

func (u *SensorRepository) sensorFieldWithoutId() string {
	return "name, unit, id_node, id_hardware"
}

func (u *SensorRepository) sensorField() string {
	return "sensor.id_sensor, sensor.name, sensor.unit, sensor.id_node, sensor.id_hardware"
}

func (u *SensorRepository) sensorPointer(sensor *entities.Sensor) []interface{} {
	return []interface{}{&sensor.IdSensor, &sensor.Name, &sensor.Unit, &sensor.IdNode, &sensor.IdHardware}
}

func (h *SensorRepository) Create(ctx context.Context, tx pgx.Tx, payload *entities.SensorCreate) (sensor entities.Sensor, err error) {
	sensor = entities.Sensor{
		IdSensor:     0,
		SensorCreate: *payload,
	}
	sqlStatement := fmt.Sprintf(`
	INSERT INTO "sensor" (
		%s
	)
	VALUES ($1, $2, $3, $4) RETURNING id_sensor`, h.sensorFieldWithoutId())
	err = tx.QueryRow(ctx, sqlStatement, sensor.Name, sensor.Unit, sensor.IdNode, sensor.IdHardware).Scan(&sensor.IdSensor)
	if err != nil {
		return sensor, err
	}

	return sensor, nil
}

func (u *SensorRepository) GetAll(ctx context.Context, tx pgx.Tx, currentUser *entities.UserRead) (sensors []entities.Sensor, err error) {
	sensors = []entities.Sensor{}
	var sqlStatement string
	var rows pgx.Rows
	if currentUser.IsAdmin {
		sqlStatement = fmt.Sprintf(`SELECT %s FROM "sensor"`, u.sensorField())
		rows, err = tx.Query(ctx, sqlStatement)
		if err != nil {
			return sensors, err
		}
		defer rows.Close()
	} else {
		sqlStatement = fmt.Sprintf(`SELECT %s FROM "sensor" INNER JOIN "node" ON node.id_node=sensor.id_node WHERE node.id_user=$1`, u.sensorField())
		rows, err = tx.Query(ctx, sqlStatement, currentUser.IdUser)
		if err != nil {
			return sensors, err
		}
		defer rows.Close()
	}

	for rows.Next() {
		var sensor entities.Sensor
		err := rows.Scan(
			u.sensorPointer(&sensor)...,
		)
		if err != nil {
			return sensors, err
		}
		sensors = append(sensors, sensor)
	}
	if err := rows.Err(); err != nil {
		return sensors, err
	}
	return sensors, nil
}

func (u *SensorRepository) GetById(ctx context.Context, tx pgx.Tx, id int) (sensor entities.Sensor, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "sensor" WHERE id_sensor=$1`, u.sensorField())
	err = tx.QueryRow(ctx, sqlStatement, id).Scan(
		u.sensorPointer(&sensor)...,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return sensor, fiber.NewError(404, fmt.Sprintf("Sensor with id %d not found", id))
		}
		return sensor, err
	}
	return sensor, nil
}

func (u *SensorRepository) GetHardwareSensor(ctx context.Context, tx pgx.Tx, hardwareId int) ([]entities.Sensor, error) {
	sensors := []entities.Sensor{}
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "sensor" WHERE id_hardware=$1`, u.sensorField())
	rows, err := tx.Query(ctx, sqlStatement, hardwareId)
	if err != nil {
		return sensors, err
	}
	defer rows.Close()

	for rows.Next() {
		var sensor entities.Sensor
		err := rows.Scan(
			u.sensorPointer(&sensor)...,
		)
		if err != nil {
			return sensors, err
		}
		sensors = append(sensors, sensor)
	}
	if err := rows.Err(); err != nil {
		return sensors, err
	}

	return sensors, nil
}

func (u *SensorRepository) GetNodeSensor(ctx context.Context, tx pgx.Tx, nodeId int) ([]entities.Sensor, error) {
	sensors := []entities.Sensor{}
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "sensor" WHERE id_node=$1`, u.sensorField())
	rows, err := tx.Query(ctx, sqlStatement, nodeId)
	if err != nil {
		return sensors, err
	}
	defer rows.Close()

	for rows.Next() {
		var sensor entities.Sensor
		err := rows.Scan(
			u.sensorPointer(&sensor)...,
		)
		if err != nil {
			return sensors, err
		}
		sensors = append(sensors, sensor)
	}
	if err := rows.Err(); err != nil {
		return sensors, err
	}

	return sensors, nil
}

func (u *SensorRepository) GetSensorChannel(ctx context.Context, tx pgx.Tx, sensorId int) (channels []entities.Channel, err error) {
	channels = []entities.Channel{}
	sqlStatement := `SELECT channel.time, channel.value, channel.id_sensor FROM "channel" WHERE channel.id_sensor=$1`
	rows, err := tx.Query(ctx, sqlStatement, sensorId)
	if err != nil {
		return channels, err
	}
	defer rows.Close()

	for rows.Next() {
		var channel entities.Channel
		err := rows.Scan(
			&channel.Time, &channel.Value, &channel.IdSensor,
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

func (u *SensorRepository) GetIdUserWhoOwnSensorById(ctx context.Context, tx pgx.Tx, sensorId int) (userId int, err error) {
	sqlStatement := `SELECT node.id_user FROM "sensor" INNER JOIN "node" ON node.id_node=sensor.id_node WHERE sensor.id_sensor=$1`
	err = tx.QueryRow(ctx, sqlStatement, sensorId).Scan(&userId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return userId, fiber.NewError(404, fmt.Sprintf("Sensor with id %d not found", sensorId))
		}
		return userId, err
	}
	return userId, nil
}

func (u *SensorRepository) Update(ctx context.Context, tx pgx.Tx, sensor *entities.Sensor, payload *entities.SensorUpdate) (err error) {
	payload.ChangeSettedFieldOnly(sensor)

	sqlStatement := `
	UPDATE "sensor"
	SET name=$1, unit=$2
	WHERE id_sensor=$3`
	res, err := tx.Exec(ctx, sqlStatement, payload.Name, payload.Unit, sensor.IdSensor)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on update sensor with id %d", sensor.IdSensor))
	}
	return nil
}

func (u *SensorRepository) Delete(ctx context.Context, tx pgx.Tx, id int) (err error) {
	sqlStatement := `DELETE FROM "sensor" WHERE id_sensor=$1`
	res, err := tx.Exec(ctx, sqlStatement, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on delete with id %d", id))
	}
	return nil
}
