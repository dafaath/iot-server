package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NodeHandler struct {
	db                 *pgxpool.Pool
	repository         *repositories.NodeRepository
	hardwareRepository *repositories.HardwareRepository
	channelRepository  *repositories.ChannelRepository
	validator          *dependencies.Validator
}

func NewNodeHandler(db *pgxpool.Pool, nodeRepository *repositories.NodeRepository, channelRepository *repositories.ChannelRepository, hardwareRepository *repositories.HardwareRepository, validator *dependencies.Validator) (NodeHandler, error) {
	return NodeHandler{
		db:                 db,
		repository:         nodeRepository,
		hardwareRepository: hardwareRepository,
		channelRepository:  channelRepository,
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

	// Make hardware validation for node async
	validateNodeHardwareChannel := make(chan error)
	go func() {
		hardwareType, err := h.hardwareRepository.GetHardwareTypeById(ctx, h.db, bodyPayload.IdHardwareNode)
		if err != nil {
			validateNodeHardwareChannel <- err
			return
		}

		hardwareType = strings.ToLower(hardwareType)
		if hardwareType != "microcontroller unit" && hardwareType != "single-board computer" {
			validateNodeHardwareChannel <- fiber.NewError(400, "Node Hardware type not match, type should be microcontroller unit or single-board computer")
			return
		}

		validateNodeHardwareChannel <- nil
	}()

	// Make hardware validation for sensor async
	withoutCurlyBrace := bodyPayload.IdHardwareSensor[1 : len(bodyPayload.IdHardwareSensor)-1]
	idStringArray := strings.Split(withoutCurlyBrace, ",")
	sensorHardwareIdLength := len(idStringArray)
	validateSensorHardwareChannel := make(chan error, sensorHardwareIdLength)
	for _, id := range idStringArray {
		id = strings.TrimSpace(id)
		go func(id string) {
			if strings.ToLower(id) == "null" {
				validateSensorHardwareChannel <- nil
				return
			}

			idInt, err := strconv.Atoi(id)
			if err != nil {
				validateSensorHardwareChannel <- fiber.NewError(400, fmt.Sprintf("Sensor Hardware id must be integer, current id is '%s'", id))
				return
			}

			hardwareType, err := h.hardwareRepository.GetHardwareTypeById(ctx, h.db, idInt)
			if err != nil {
				validateSensorHardwareChannel <- err
				return
			}
			if hardwareType != "sensor" {
				validateSensorHardwareChannel <- fiber.NewError(400, fmt.Sprintf("Sensor Hardware type for id %s not match, type should be sensor", id))
				return
			}

			validateSensorHardwareChannel <- nil
		}(id)
	}

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	err = <-validateNodeHardwareChannel
	if err != nil {
		return err
	}

	for i := 0; i < sensorHardwareIdLength; i++ {
		err = <-validateSensorHardwareChannel
		if err != nil {
			return err
		}
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

	// Get all channel for each node
	limit, err := strconv.Atoi(c.Query("limit", "-1"))
	if err != nil {
		return fiber.NewError(400, "Limit must be integer")
	}

	type NodeWithChannelOutput struct {
		NodeWithChannel entities.NodeWithChannel
		Err             error
	}
	nodeWithChannelChannel := make(chan NodeWithChannelOutput, len(nodes))
	for _, node := range nodes {
		go func(node entities.Node) {
			channel, err := h.channelRepository.GetNodeChannel(ctx, h.db, node.IdNode, limit)
			nodeWithChannel := entities.NodeWithChannel{
				Node: node,
				Feed: channel,
			}
			nodeWithChannelChannel <- NodeWithChannelOutput{
				NodeWithChannel: nodeWithChannel,
				Err:             err,
			}
		}(node)
	}

	nodesWithChannel := make([]entities.NodeWithChannel, len(nodes))
	for i := 0; i < len(nodes); i++ {
		nodeWithChannel := <-nodeWithChannelChannel
		if nodeWithChannel.Err != nil {
			return nodeWithChannel.Err
		}

		nodesWithChannel[i] = nodeWithChannel.NodeWithChannel
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("node", fiber.Map{
			"nodes": nodesWithChannel,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(nodesWithChannel)
	}
}

func (h *NodeHandler) GetById(c *fiber.Ctx) (err error) {
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

	type NodeResponse struct {
		node entities.Node
		err  error
	}
	nodeResponseChannel := make(chan NodeResponse)
	go func() {
		node, err := h.repository.GetById(ctx, tx, id)
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

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin && !node.IsPublic {
		return fiber.NewError(403, "You can't see another user's node")
	}

	limit, err := strconv.Atoi(c.Query("limit", "-1"))
	if err != nil {
		return fiber.NewError(400, "Limit must be integer")
	}

	feed, err := h.channelRepository.GetNodeChannel(ctx, tx, node.IdNode, limit)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		return c.Render("node_detail", fiber.Map{
			"node": node,
			"feed": feed,
		}, "layouts/main")
	default:
		return c.Status(fiber.StatusOK).JSON(entities.NodeWithChannel{
			Node: node,
			Feed: feed,
		})
	}
}

func (h *NodeHandler) UpdateForm(c *fiber.Ctx) (err error) {
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

	node, err := h.repository.GetById(ctx, tx, id)
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

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	node, err := h.repository.GetById(ctx, tx, id)
	if err != nil {
		return err
	}

	currentUser, err := h.validator.GetAuthentication(c)
	if err != nil {
		return err
	}

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin {
		return fiber.NewError(403, "Can't edit another user's data")
	}

	err = h.repository.Update(ctx, tx, &node, bodyPayload)
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

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(ctx, tx, &err)

	node, err := h.repository.GetById(ctx, tx, id)
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

	err = h.repository.Delete(ctx, tx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Success delete node, id: %d", id))
}
