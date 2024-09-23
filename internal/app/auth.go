package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/recommender-system-for-MTUCI/backend/internal/models"
)

func (ctrl *Controller) handleRegistration(c echo.Context) error {
	var request models.RegistrationRequst
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"erroe": err.Error()})
	}
	return c.JSON(http.StatusCreated, request)
}
