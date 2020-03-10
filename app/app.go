package app

import (
	"encoding/gob"
	"time"

	"github.com/jinzhu/gorm"

	// SQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"

	"github.com/gorilla/sessions"
	"github.com/iammarkps/tu83-announcer-api/handler"
	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
)

// New function create new echo app
func New() (*echo.Echo, *gorm.DB) {
	gob.Register(time.Time{})
	e := echo.New()

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=announcer sslmode=disable")
	e.Logger.Info("Connecting to database...")

	if err != nil {
		panic("Failed to connect to database")
	}

	e.Logger.Info("Successfully connected to database")
	// defer db.Close()
	db.AutoMigrate(&models.User{})

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowOrigins:     []string{"http://localhost:3000"},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	e.Debug = true

	h := &handler.Handler{DB: db}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "😼Triam Udom Suksa School's exam announcer API is running!")
	})

	e.GET("/student", h.Student)
	e.POST("/confirm", h.Confirm)
	e.POST("/login", h.Login)
	e.GET("/logout", h.Logout)

	return e, db
}
