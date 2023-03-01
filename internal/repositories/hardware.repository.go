package repositories

import (
	"context"
	"fmt"

	"github.com/dafaath/iot-server/internal/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type HardwareRepository struct{}

func NewHardwareRepository() (HardwareRepository, error) {
	return HardwareRepository{}, nil
}

func (u *HardwareRepository) hardwareField() string {
	return "id_hardware, name, type, description"
}

func (u *HardwareRepository) hardwarePointer(hardware *entities.Hardware) []interface{} {
	return []interface{}{&hardware.IdHardware, &hardware.Name, &hardware.Type, &hardware.Description}
}

func (h *HardwareRepository) Create(ctx context.Context, tx pgx.Tx, payload *entities.HardwareCreate) (hardware entities.Hardware, err error) {
	hardware = entities.Hardware{
		HardwareCreate: *payload,
	}
	sqlStatement := `
	INSERT INTO "hardware" (
		name,
		type,
		description
	)
	VALUES ($1, $2, $3) RETURNING id_hardware`
	err = tx.QueryRow(ctx, sqlStatement, hardware.Name, hardware.Type, hardware.Description).Scan(&hardware.IdHardware)
	if err != nil {
		return hardware, err
	}

	return hardware, nil
}

func (u *HardwareRepository) getAllItem(ctx context.Context, tx pgx.Tx, sqlStatement string) (hardwares []entities.Hardware, err error) {
	hardwares = []entities.Hardware{}
	rows, err := tx.Query(ctx, sqlStatement)
	if err != nil {
		return hardwares, err
	}
	defer rows.Close()

	for rows.Next() {
		var hardware entities.Hardware
		err := rows.Scan(
			u.hardwarePointer(&hardware)...,
		)
		if err != nil {
			return hardwares, err
		}
		hardwares = append(hardwares, hardware)
	}
	if err := rows.Err(); err != nil {
		return hardwares, err
	}
	return hardwares, nil

}

func (u *HardwareRepository) GetAllHardware(ctx context.Context, tx pgx.Tx) (hardwares []entities.Hardware, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "hardware"`, u.hardwareField())
	return u.getAllItem(ctx, tx, sqlStatement)
}

func (u *HardwareRepository) GetAllNode(ctx context.Context, tx pgx.Tx) (hardwares []entities.Hardware, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "hardware" WHERE lower(type) = 'single-board computer' or lower(type) = 'microcontroller unit'`, u.hardwareField())
	return u.getAllItem(ctx, tx, sqlStatement)
}
func (u *HardwareRepository) GetAllSensor(ctx context.Context, tx pgx.Tx) (hardwares []entities.Hardware, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "hardware" WHERE lower(type) = 'sensor'`, u.hardwareField())
	return u.getAllItem(ctx, tx, sqlStatement)
}

func (u *HardwareRepository) GetById(ctx context.Context, tx pgx.Tx, id int) (hardware entities.Hardware, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "hardware" WHERE id_hardware=$1`, u.hardwareField())
	err = tx.QueryRow(ctx, sqlStatement, id).Scan(
		u.hardwarePointer(&hardware)...,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return hardware, fiber.NewError(404, fmt.Sprintf("Hardware with id %d not found", id))
		}
		return hardware, err
	}
	return hardware, nil
}

func (u *HardwareRepository) Update(ctx context.Context, tx pgx.Tx, hardware *entities.Hardware, payload *entities.HardwareUpdate) (err error) {
	payload.ChangeSettedFieldOnly(hardware)

	sqlStatement := `
	UPDATE "hardware"
	SET name=$1, type=$2, description=$3
	WHERE id_hardware=$4`
	res, err := tx.Exec(ctx, sqlStatement, payload.Name, payload.Type, payload.Description, hardware.IdHardware)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on update hardware with id %d", hardware.IdHardware))
	}
	return nil
}

func (u *HardwareRepository) Delete(ctx context.Context, tx pgx.Tx, id int) (err error) {
	sqlStatement := `DELETE FROM "hardware" WHERE id_hardware=$1`
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
