package handlers

import (
	"context"

	"github.com/dafaath/iot-server/internal/dependencies"
	"github.com/dafaath/iot-server/internal/entities"
	"github.com/dafaath/iot-server/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChannelHandler struct {
	db               *pgxpool.Pool
	repository       *repositories.ChannelRepository
	sensorRepository *repositories.SensorRepository
	validator        *dependencies.Validator
}

func NewChannelHandler(db *pgxpool.Pool, channelRepository *repositories.ChannelRepository, sensorRepository *repositories.SensorRepository, validator *dependencies.Validator) (ChannelHandler, error) {
	return ChannelHandler{
		db:               db,
		repository:       channelRepository,
		sensorRepository: sensorRepository,
		validator:        validator,
	}, nil
}

func (h *ChannelHandler) CreateForm(c *fiber.Ctx) (err error) {
	idSensor := c.QueryInt("id_sensor", 0)
	return c.Render("channel_form", fiber.Map{"title": "Create Channel", "idSensor": idSensor}, "layouts/main")
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

	// Get sensor owner Id async
	type sensorOwnerIdResult struct {
		res int
		err error
	}
	sensorOwnerIdChannel := make(chan sensorOwnerIdResult)
	go func() {
		sensorOwnerId, err := h.sensorRepository.GetIdUserWhoOwnSensorById(ctx, h.db, bodyPayload.IdSensor)
		sensorOwnerIdChannel <- sensorOwnerIdResult{
			res: sensorOwnerId,
			err: err,
		}
	}()

	currentUserRes := <-currentUserChannel
	currentUser := currentUserRes.res
	err = currentUserRes.err
	if err != nil {
		return err
	}

	sensorOwnerIdRes := <-sensorOwnerIdChannel
	sensorOwnerId := sensorOwnerIdRes.res
	err = sensorOwnerIdRes.err
	if err != nil {
		return err
	}

	if currentUser.IdUser != sensorOwnerId {
		return fiber.NewError(fiber.StatusForbidden, "You can't send channel to another user's sensor")
	}

	_, err = h.repository.Create(ctx, h.db, &bodyPayload)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).SendString("Add new channel")

}
