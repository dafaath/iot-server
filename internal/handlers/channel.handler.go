package handlers

import (
	"context"

	"github.com/dafaath/iot-server/v2/internal/dependencies"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChannelHandler struct {
	db             *pgxpool.Pool
	repository     *repositories.ChannelRepository
	nodeRepository *repositories.NodeRepository
	validator      *dependencies.Validator
}

func NewChannelHandler(db *pgxpool.Pool, channelRepository *repositories.ChannelRepository, nodeRepository *repositories.NodeRepository, validator *dependencies.Validator) (ChannelHandler, error) {
	return ChannelHandler{
		db:             db,
		repository:     channelRepository,
		nodeRepository: nodeRepository,
		validator:      validator,
	}, nil
}

func (h *ChannelHandler) CreateForm(c *fiber.Ctx) (err error) {
	idNode := c.QueryInt("id_node", 0)
	return c.Render("channel_form", fiber.Map{"title": "Create Channel", "idNode": idNode}, "layouts/main")
}

func (h *ChannelHandler) Create(c *fiber.Ctx) (err error) {
	ctx := context.Background()
	bodyPayload := entities.ChannelCreate{}

	parseChannel := make(chan error)
	go func() {
		err = h.validator.ParseBody(c, &bodyPayload)
		parseChannel <- err
	}()

	// Get current user async
	type currentUserResult struct {
		res entities.UserRead
		err error
	}
	currentUserChannel := make(chan currentUserResult)
	go func() {
		currentUser, err := h.validator.GetAuthentication(c)
		currentUserChannel <- currentUserResult{
			res: currentUser,
			err: err,
		}
	}()

	// Wait for parsing and get parsing error
	err = <-parseChannel
	if err != nil {
		return err
	}

	// Validate node owner async
	nodeOwnerErrorChannel := make(chan error)
	go func() {
		node, err := h.nodeRepository.GetById(ctx, h.db, bodyPayload.IdNode)
		if err != nil {
			nodeOwnerErrorChannel <- err
			return
		}

		currentUserRes := <-currentUserChannel
		err = currentUserRes.err
		if err != nil {
			nodeOwnerErrorChannel <- err
			return
		}
		currentUser := currentUserRes.res

		if currentUser.IdUser != node.IdUser {
			nodeOwnerErrorChannel <- fiber.NewError(fiber.StatusForbidden, "You can't send channel to another user's node")
			return
		}

		nodeOwnerErrorChannel <- nil
	}()

	err = <-nodeOwnerErrorChannel
	if err != nil {
		return err
	}

	err = h.repository.Create(ctx, h.db, &bodyPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString("Add new channel")

}
