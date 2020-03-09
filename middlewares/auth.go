package middlewares

import (
	"fmt"

	"log"

	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Auth middleware
func (middleware *Middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("SESSION", c)
		User := &models.User{}

		ID := fmt.Sprintf("%v", sess.Values["user"])
		log.Print(sess.Values)
		middleware.DB.Where(&models.User{ID: ID}).First(User)

		if User != (&models.User{}) {
			c.Set("User", User)
			return next(c)
		}

		return echo.ErrUnauthorized
	}
}
