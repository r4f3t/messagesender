package controller

import (
	"github.com/labstack/echo/v4"
)

func MakeHandler(instance *echo.Echo, s *resource) {
	e := instance

	e.POST("/start", s.StartProcessingMessages)
	e.POST("/stop", s.StopProcessingMessages)

	e.GET("/messages/sent", s.GetSentMessages)

}
