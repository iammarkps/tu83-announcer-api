package app

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	// Redis
	_ "github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"gopkg.in/boj/redistore.v1"

	// SQL Dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/iammarkps/tu83-announcer-api/handler"
	"github.com/iammarkps/tu83-announcer-api/models"
	"github.com/labstack/echo-contrib/session"
)

// New function create new echo app
func New() (*echo.Echo, *gorm.DB) {
	gob.Register(time.Time{})
	e := echo.New()

	db, err := gorm.Open("postgres", os.Getenv("DB"))
	e.Logger.Info("Connecting to database...")

	if err != nil {
		panic(err)
	}

	e.Logger.Info("Successfully connected to database")
	db.BlockGlobalUpdate(true)
	db.DB().SetMaxIdleConns(100)

	// defer db.Close()
	db.AutoMigrate(&models.User{})

	e.Use(session.Middleware(newRedisStore()))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowOrigins:     []string{"http://localhost:3000", "https://announce.triamudom.ac.th", "https://triam83announce.netlify.com"},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())

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

func newRedisStore() sessions.Store {
	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte(os.Getenv("REDIS_KEY")))

	if err != nil {
		panic(err)
	} else {
		log.Printf("Successfully connected to redis")
	}

	return store
}
