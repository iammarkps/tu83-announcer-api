package handler

import (
	"fmt"
	"net/http"

	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Student handler
func (handler *Handler) Student(c echo.Context) error {
	sess, _ := session.Get("Session", c)
	User := &models.User{}

	ID := fmt.Sprintf("%v", sess.Values["user"])
	handler.DB.Where(&models.User{ID: ID}).First(User)

	if User == (&models.User{}) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	return c.JSON(http.StatusOK, User)
}
