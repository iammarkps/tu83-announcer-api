package handler

import (
	"fmt"
	"net/http"

	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Confirm handler
func (handler *Handler) Confirm(c echo.Context) error {
	sess, _ := session.Get("Session", c)
	User := &models.User{}

	ID := fmt.Sprintf("%v", sess.Values["user"])
	handler.DB.Where(&models.User{ID: ID}).First(User)

	if User.Confirmed {
		return c.JSON(http.StatusUnauthorized, "Already")
	}

	handler.DB.Model(User).Update("Confirmed", true)
	return c.JSON(http.StatusOK, "OK")
}
