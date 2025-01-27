package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/r4f3t/messagesender/internal/message"
)

type resource struct {
	service message.MessageService
}

func NewController(service message.MessageService) *resource {
	return &resource{
		service: service,
	}
}

// ProcessingMessage godoc
// @Summary Starts Message Processings
// @Description No return value
// @Tags Message Processing
// @Accept json
// @Produce json
// @Success 200 {string} string "Message processing started."
// @Router /start [post]
func (receiver *resource) StartProcessingMessages(c echo.Context) error {
	err := receiver.service.StartProcessingMessages()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Message processing started.",
	})
}

// ProcessingMessage godoc
// @Summary Stops Message Processings
// @Description No return value
// @Tags Message Processing
// @Accept json
// @Produce json
// @Success 200 {string} string "Message processing stopped."
// @Router /stop [post]
func (receiver *resource) StopProcessingMessages(c echo.Context) error {
	err := receiver.service.StopProcessingMessages()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Message processing stopped.",
	})
}

// ProcessingMessage godoc
// @Summary Retrive list of sent messages
// @Description Returns list of sent messages
// @Tags Message Processing
// @Accept json
// @Produce json
// @Success 200 {array} message.Message "List of sent messages"
// @Router /messages/sent [get]
func (receiver *resource) GetSentMessages(c echo.Context) error {
	messages, err := receiver.service.GetSentMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if len(messages) == 0 {
		return c.JSON(http.StatusNotFound, map[string][]message.Message{
			"message": messages,
		})
	}

	return c.JSON(http.StatusOK, map[string][]message.Message{
		"message": messages,
	})
}
