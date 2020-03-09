package app

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"

	// SQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"

	"github.com/gorilla/sessions"
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

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "😼Triam Udom Suksa School's exam announcer API is running!")
	})

	e.GET("/student", func(c echo.Context) error {
		sess, _ := session.Get("SESSION", c)
		User := &models.User{}

		ID := fmt.Sprintf("%v", sess.Values["user"])
		db.Where(&models.User{ID: ID}).First(User)

		return c.JSON(http.StatusOK, User)
	})

	e.POST("/login", func(c echo.Context) error {
		type userRes struct {
			ID string `json:"id"`
		}

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

		sess, _ := session.Get("SESSION", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
		}

		sess.Values["user"] = User.ID
		sess.Values["timestamp"] = time.Now()

		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			e.Logger.Fatal(err)
			return c.JSON(http.StatusInternalServerError, "Cannot logged you in")
		}

		return c.JSON(http.StatusOK, &userRes{ID: User.ID})
	})

	e.GET("/logout", func(c echo.Context) error {
		sess, _ := session.Get("SESSION", c)
		sess.Options.MaxAge = -1

		err := sess.Save(c.Request(), c.Response())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Cannot logged you out")
		}

		return c.JSON(http.StatusOK, "OK")
	})

	return e, db
}
