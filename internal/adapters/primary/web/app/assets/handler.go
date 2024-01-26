package assets

import (
	"github.com/labstack/echo/v4"
)

const (
	// Routes

	AppletRoute = "/assets"
)

type Handler interface {
	Applet(c echo.Context) error
}

type handler struct {
}

func New() Handler {
	return &handler{}
}
