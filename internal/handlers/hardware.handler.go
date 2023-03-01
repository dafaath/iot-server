package handlers

import (
	"context"
	"fmt"

	"github.com/dafaath/iot-server/internal/dependencies"
	"github.com/dafaath/iot-server/internal/entities"
	"github.com/dafaath/iot-server/internal/helper"
	"github.com/dafaath/iot-server/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HardwareHandler struct {
	db               *pgxpool.Pool
	repository       *repositories.HardwareRepository
	validator        *dependencies.Validator
	nodeRepository   *repositories.NodeRepository
	sensorRepository *repositories.SensorRepository
}

func NewHardwareHandler(db *pgxpool.Pool, hardwareRepository *repositories.HardwareRepository, nodeRepository *repositories.NodeRepository, sensorRepository *repositories.SensorRepository, validator *dependencies.Validator) (HardwareHandler, error) {
	return HardwareHandler{
		db:               db,
		validator:        validator,
		repository:       hardwareRepository,
		nodeRepository:   nodeRepository,
		sensorRepository: sensorRepository,
	}, nil
}

func (h *HardwareHandler) CreateForm(c *fiber.Ctx) (err error) {
	return c.Render("hardware_form", fiber.Map{
		"title": "Add Hardware",
	}, "layouts/main")
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

	nodes, err := h.repository.GetAllNode(ctx, tx)
	if err != nil {
		return err
	}

	sensors, err := h.repository.GetAllSensor(ctx, tx)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("hardware", fiber.Map{
			"title":  "Hardware",
			"node":   nodes,
			"sensor": sensors,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"node":   nodes,
			"sensor": sensors,
		})

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

	switch hardware.Type {
	case "microcontroller unit", "single-board computer":
		nodes, err := h.nodeRepository.GetHardwareNode(ctx, tx, hardware.IdHardware)
		if err != nil {
			return err
		}

		accept := c.Accepts("application/json", "text/html")
		switch accept {
		case "text/html":
			return c.Render("hardware_detail", fiber.Map{
				"title":    "Hardware Detail",
				"hardware": hardware,
				"nodes":    nodes,
			}, "layouts/main")
		default:
			return c.Status(fiber.StatusOK).JSON(entities.HardwareWithNode{
				Hardware: hardware,
				Nodes:    nodes,
			})
		}

	case "sensor":
		sensors, err := h.sensorRepository.GetHardwareSensor(ctx, tx, id)
		if err != nil {
			return err
		}
		accept := c.Accepts("application/json", "text/html")
		switch accept {
		case "text/html":
			return c.Render("hardware_detail", fiber.Map{
				"title":    "Hardware Detail",
				"hardware": hardware,
				"sensors":  sensors,
			}, "layouts/main")
		default:
			return c.Status(fiber.StatusOK).JSON(entities.HardwareWithSensor{
				Hardware: hardware,
				Sensors:  sensors,
			})
		}

	default:
		return c.Status(fiber.StatusOK).JSON(hardware)
	}

}

func (h *HardwareHandler) UpdateForm(c *fiber.Ctx) (err error) {
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}
	ctx := context.Background()

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	hardware, err := h.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	return c.Render("hardware_form", fiber.Map{
		"title":    "Update Hardware",
		"hardware": hardware,
		"edit":     true,
	}, "layouts/main")
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
