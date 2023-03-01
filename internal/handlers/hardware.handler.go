package handlers

import (
	"context"
	"fmt"

	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HardwareHandler struct {
	db             *pgxpool.Pool
	repository     *repositories.HardwareRepository
	validator      *dependencies.Validator
	nodeRepository *repositories.NodeRepository
}

func NewHardwareHandler(db *pgxpool.Pool, hardwareRepository *repositories.HardwareRepository, nodeRepository *repositories.NodeRepository, validator *dependencies.Validator) (HardwareHandler, error) {
	return HardwareHandler{
		db:             db,
		validator:      validator,
		repository:     hardwareRepository,
		nodeRepository: nodeRepository,
	}, nil
}

func (h *HardwareHandler) Create(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	bodyPayload := &entities.HardwareCreate{}

	err = h.validator.ParseBody(c, bodyPayload)
	if err != nil {
		return err
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	_, err = h.repository.Create(ctx, tx, bodyPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString("Success add new hardware")
}

func (h *HardwareHandler) GetAll(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	hardwares, err := h.repository.GetAllHardware(ctx, tx)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("hardware", fiber.Map{
			"hardwares": hardwares,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(hardwares)
	}

}

func (h *HardwareHandler) GetById(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	hardware, err := h.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("hardware_detail", fiber.Map{
			"hardware": hardware,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(hardware)
	}

}

func (h *HardwareHandler) Update(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	bodyPayload := &entities.HardwareUpdate{}
	err = h.validator.ParseBody(c, bodyPayload)
	if err != nil {
		return err
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	hardware, err := h.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = h.repository.Update(ctx, tx, &hardware, bodyPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("Success edit hardware")
}

func (h *HardwareHandler) Delete(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	_, err = h.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	err = h.repository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Success delete hardware, id: %d", id))
}
