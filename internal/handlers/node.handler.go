package handlers

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dafaath/iot-server/internal/dependencies"
	"github.com/dafaath/iot-server/internal/entities"
	"github.com/dafaath/iot-server/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeHandler struct {
	db                 *pgxpool.Pool
	repository         *repositories.NodeRepository
	hardwareRepository *repositories.HardwareRepository
	sensorRepository   *repositories.SensorRepository
	validator          *dependencies.Validator
}

func NewNodeHandler(db *pgxpool.Pool, nodeRepository *repositories.NodeRepository, hardwareRepository *repositories.HardwareRepository, sensorRepository *repositories.SensorRepository, validator *dependencies.Validator) (NodeHandler, error) {
	return NodeHandler{
		db:                 db,
		repository:         nodeRepository,
		hardwareRepository: hardwareRepository,
		sensorRepository:   sensorRepository,
		validator:          validator,
	}, nil
}

func (h *NodeHandler) CreateForm(c *fiber.Ctx) (err error) {
	return c.Render("node_form", fiber.Map{
		"title": "Create Node",
	}, "layouts/main")
}

func (h *NodeHandler) Create(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	bodyPayload := entities.NodeCreate{}
	parseChannel := make(chan error)

	// Make parse body async in background
	go func() {
		err = h.validator.ParseBody(c, &bodyPayload)
		parseChannel <- err
	}()

	err = <-parseChannel
	if err != nil {
		return err
	}

	// Make hardware validation async in background
	validateHardwareChannel := make(chan error)
	go func() {
		hardware, err := h.hardwareRepository.GetById(ctx, h.db, bodyPayload.IdHardware)
		if err != nil {
			validateHardwareChannel <- err
			return
		}

		hardwareType := strings.ToLower(hardware.Type)
		if hardwareType != "microcontroller unit" && hardwareType != "single-board computer" {
			validateHardwareChannel <- fiber.NewError(400, "Hardware type not match, type should be microcontroller unit or single-board computer")
			return
		}

		validateHardwareChannel <- nil
	}()

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	err = <-validateHardwareChannel
	if err != nil {
		return err
	}

	_, err = h.repository.Create(ctx, h.db, &bodyPayload, &currentUser)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString("Success add new node")
}

func (h *NodeHandler) GetAll(c *fiber.Ctx) (err error) {
	ctx := context.Background()

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	nodes, err := h.repository.GetAll(ctx, h.db, &currentUser)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})

		return c.Render("node", fiber.Map{
			"title": "Node",
			"nodes": nodes,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(nodes)
	}
}

func (h *NodeHandler) GetById(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	type NodeResponse struct {
		node entities.Node
		err  error
	}
	nodeResponseChannel := make(chan NodeResponse)
	go func() {
		node, err := h.repository.GetById(ctx, h.db, id)
		nodeResponseChannel <- NodeResponse{node, err}
	}()

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	nodeResponse := <-nodeResponseChannel
	node := nodeResponse.node
	err = nodeResponse.err
	if err != nil {
		return err
	}

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin {
		return fiber.NewError(403, "You can’t see another user’s node")
	}

	hardware, err := h.hardwareRepository.GetById(ctx, h.db, node.IdHardware)
	if err != nil {
		return err
	}

	sensors, err := h.sensorRepository.GetNodeSensor(ctx, h.db, node.IdNode)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("node_detail", fiber.Map{
			"title":    "Node Detail",
			"node":     node,
			"hardware": hardware,
			"sensor":   sensors,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(entities.NodeWithHardwareAndSensors{
			Node:     node,
			Hardware: hardware,
			Sensor:   sensors,
		})
	}
}

func (h *NodeHandler) UpdateForm(c *fiber.Ctx) (err error) {
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}
	ctx := context.Background()

	node, err := h.repository.GetById(ctx, h.db, id)
	if err != nil {
		return err
	}

	return c.Render("node_form", fiber.Map{
		"title": "Edit Node",
		"node":  node,
		"edit":  true,
	}, "layouts/main")
}

func (h *NodeHandler) Update(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	bodyPayload := &entities.NodeUpdate{}
	err = h.validator.ParseBody(c, bodyPayload)
	if err != nil {
		return err
	}

	node, err := h.repository.GetById(ctx, h.db, id)
	if err != nil {
		return err
	}

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin {
		return fiber.NewError(403, "Can’t edit another user’s data")
	}

	err = h.repository.Update(ctx, h.db, &node, bodyPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString("Success edit node")
}

func (h *NodeHandler) Delete(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	id, err := h.validator.ParseIdFromUrlParameter(c)
	if err != nil {
		return err
	}

	node, err := h.repository.GetById(ctx, h.db, id)
	if err != nil {
		return err
	}

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin {
		return fiber.NewError(403, "You can’t delete another user’s node")
	}

	err = h.repository.Delete(ctx, h.db, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Success delete node, id: %d", id))
}
