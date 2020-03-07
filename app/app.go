package app

import (
	"github.com/jinzhu/gorm"

	// SQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"

	"github.com/iammarkps/tu83-announcer-api/models"
)

// New function create new echo app
func New() (*echo.Echo, *gorm.DB) {
	e := echo.New()

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=announcer sslmode=disable")
	e.Logger.Info("Connecting to database...")

	if err != nil {
		panic("Failed to connect to database")
	}

	e.Logger.Info("Successfully connected to database")

	// defer db.Close()

	db.AutoMigrate(&models.User{})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

	e.Debug = true

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "😼Triam Udom Suksa School's exam announcer API is running!")
	})

	e.POST("/student", func(c echo.Context) error {
		type userReq struct {
			ID    string `json:"id"`
			CtzID string `json:"ctz_id"`
		}

		u := new(userReq)

		if err := c.Bind(&u); err != nil {
			return err
		}

		User := &models.User{}

		db.Where(&models.User{ID: u.ID}).First(User)

		if !(User.CtzID == u.CtzID) || User == (&models.User{}) {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}

		return c.JSON(http.StatusOK, User)
	})

	return e, db
}
