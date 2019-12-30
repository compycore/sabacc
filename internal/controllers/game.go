package controllers

import (
	"github.com/labstack/echo"
)

func Play(c echo.Context) error {
	// TODO Parse URI

	// TODO Do I need/want to send a response back?
	return c.JSON(200, nil)
}
