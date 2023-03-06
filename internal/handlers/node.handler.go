package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
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
	ctx := context.Background()
	nodeHardware, err := h.hardwareRepository.GetAllNode(ctx, h.db)
	if err != nil {
		return err
	}

	sensorHardware, err := h.hardwareRepository.GetAllSensor(ctx, h.db)
	if err != nil {
		return err
	}

	return c.Render("node_form", fiber.Map{
		"title":          "Create Node",
		"nodeHardware":   nodeHardware,
		"sensorHardware": sensorHardware,
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

	// Validate sensor hardware id same length with sensor field
	if len(bodyPayload.IdHardwareSensor) != len(bodyPayload.FieldSensor) {
		return fiber.NewError(400, "Sensor Hardware Id and Sensor Field must have same length")
	}

	// Make hardware validation for sensor async
	sensorHardwareIdLength := len(bodyPayload.IdHardwareSensor)
	validateSensorHardwareChannel := make(chan error, sensorHardwareIdLength)
	for _, id := range bodyPayload.IdHardwareSensor {
		go func(id int) {
			hardwareType, err := h.hardwareRepository.GetHardwareTypeById(ctx, h.db, id)
			if err != nil {
				validateSensorHardwareChannel <- err
				return
			}
			if hardwareType != "sensor" {
				validateSensorHardwareChannel <- fiber.NewError(400, fmt.Sprintf("Sensor Hardware type for id %d not match, type should be sensor", id))
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
		sort.Slice(nodesWithChannel, func(i, j int) bool {
			return nodesWithChannel[i].Node.Name < nodesWithChannel[j].Node.Name
		})
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

	if node.IdUser != currentUser.IdUser && !currentUser.IsAdmin && !node.IsPublic {
		return fiber.NewError(403, "You can't see another user's node")
	}

	limit, err := strconv.Atoi(c.Query("limit", "-1"))
	if err != nil {
		return fiber.NewError(400, "Limit must be integer")
	}

	feed, err := h.channelRepository.GetNodeChannel(ctx, h.db, node.IdNode, limit)
	if err != nil {
		return err
	}

	accept := c.Accepts("application/json", "text/html")
	switch accept {
	case "text/html":
		sensor := []fiber.Map{}
		for i := 0; i < len(node.IdHardwareSensor); i++ {
			field := ""
			if i < len(node.FieldSensor) {
				field = node.FieldSensor[i]
			}

			sensor = append(sensor, fiber.Map{
				"idHardware": node.IdHardwareSensor[i],
				"field":      field,
			})
		}

		sort.Slice(feed, func(i, j int) bool {
			return feed[i].Time.Before(feed[j].Time)
		})

		mappedChannel := []fiber.Map{}
		for i := 0; i < 10; i++ {
			mappedChannel = append(mappedChannel, fiber.Map{
				"name": fmt.Sprintf("Sensor %d", i+1),
				"data": []interface{}{},
			})
		}

		for _, channel := range feed {
			// Convert time to epoch milliseconds
			for i := 0; i < len(channel.Value); i++ {
				value := channel.Value[i]
				mappedChannel[i]["data"] = append(mappedChannel[i]["data"].([]interface{}), []interface{}{
					channel.Time.UnixMilli(),
					value,
				})
			}
		}

		channelJSONByte, err := json.Marshal(mappedChannel)
		if err != nil {
			return err
		}

		channelJSONString := string(channelJSONByte)

		return c.Render("node_detail", fiber.Map{
			"node":   node,
			"sensor": sensor,
			"feed":   channelJSONString,
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

	nodeHardware, err := h.hardwareRepository.GetAllNode(ctx, h.db)
	if err != nil {
		return err
	}

	sensorHardware, err := h.hardwareRepository.GetAllSensor(ctx, h.db)
	if err != nil {
		return err
	}

	node, err := h.repository.GetById(ctx, h.db, id)
	if err != nil {
		return err
	}

	return c.Render("node_form", fiber.Map{
		"title":          "Edit Node",
		"node":           node,
		"nodeHardware":   nodeHardware,
		"sensorHardware": sensorHardware,
		"edit":           true,
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
		return fiber.NewError(403, "Can't edit another user's data")
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
