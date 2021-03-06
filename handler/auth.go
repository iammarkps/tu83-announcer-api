package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Login handler
func (handler *Handler) Login(c echo.Context) error {
	type userRes struct {
		ID string `json:"id"`
	}

	type userReq struct {
		CtzID string `json:"ctz_id"`
		Phone string `json:"phone"`
	}

	u := new(userReq)

	if err := c.Bind(&u); err != nil {
		return err
	}

	User := &models.User{}

	handler.DB.Where(&models.User{CtzID: u.CtzID}).First(User)

	if !(User.Phone == u.Phone) || User == (&models.User{}) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	sess, _ := session.Get("Session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		SameSite: 2,
		Secure:   true,
		HttpOnly: true,
	}

	sess.Values["user"] = User.ID
	sess.Values["timestamp"] = time.Now()

	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot logged you in")
	}

	return c.JSON(http.StatusOK, &userRes{ID: User.ID})
}

// Logout handler
func (handler *Handler) Logout(c echo.Context) error {
	sess, _ := session.Get("Session", c)
	sess.Options.MaxAge = -1

	err := sess.Save(c.Request(), c.Response())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Cannot logged you out")
	}

	return c.JSON(http.StatusOK, "OK")
}
