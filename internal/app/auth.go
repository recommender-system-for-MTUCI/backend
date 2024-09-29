package app

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/recommender-system-for-MTUCI/backend/internal/models"
	"go.uber.org/zap"
)

func (ctrl *Controller) handleRegistration(c echo.Context) error {
	var request models.RegistrationRequst
	if err := c.Bind(&request); err != nil {
		ctrl.log.Error("failed to bind request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	ctrl.log.Info("bind request successfull")
	//err := ctrl.repo.User().GetUserEmail(c.Request().Context(), request.Email)
	//if err == nil {
	//	ctrl.log.Error("user with this email already in storage", zap.Error(err))
	//	return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
	//}
	hashPassword, err := hashpassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user := &models.RegistrationDTO{
		ID:       uuid.New(),
		Email:    request.Email,
		Password: string(hashPassword),
		Active:   false,
	}
	//err = ctrl.repo.User().CreateUser(c.Request().Context(), user.ID, user.Email, user.Password, user.Active)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	accsessToken, err := ctrl.tokenProvider.CreateTokenForUser(user.ID, true)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	refreshToken, err := ctrl.tokenProvider.CreateTokenForUser(user.ID, false)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	response := models.RegistrationResponse{
		ID:           user.ID,
		AccessToken:  accsessToken,
		RefreshToken: refreshToken,
	}
	return c.JSON(http.StatusCreated, response)
}

func (ctrl *Controller) handleAcceptEmail(c echo.Context) error {
	//need add logic for check token
	var request models.EmailRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, zap.Error(err))
	}
	return nil
}
