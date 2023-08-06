package httpserver

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/guneyin/bist-tools/internal/handler"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/pkg/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST              = "POST"
	PATCH             = "PATCH"
	DELETE            = "DELETE"
)

type server struct {
	app  *fiber.App
	api  fiber.Router
	port string
}

var srv *server

func Init() error {
	f := fiber.New(fiber.Config{
		ServerHeader:      "Fiber",
		AppName:           "BIST Tools",
		EnablePrintRoutes: false,
		ReadTimeout:       _defaultReadTimeout,
		WriteTimeout:      _defaultWriteTimeout,
	})

	f.Use(cors.New())

	srv = &server{
		app:  f,
		api:  f.Group("/api"),
		port: config.Cfg.HTTP.Port,
	}

	srv.setRoutes()

	return srv.start()
}

func (s *server) setRoutes() {
	// General
	s.registerHandler("/", GET, handler.GeneralStatus)

	// Auth
	s.registerHandler("/auth/login", POST, handler.Login)

	// User
	s.registerHandler("/user/me", GET, middleware.Protected(), handler.UserMe)
	s.registerHandler("/user/:id", GET, handler.GetUser)
	s.registerHandler("/user", POST, handler.CreateUser)
	s.registerHandler("/user/:id", PATCH, middleware.Protected(), handler.UpdateUser)
	s.registerHandler("/user/:id", DELETE, middleware.Protected(), handler.DeleteUser)

	// Broker
	s.registerHandler("/broker/list", GET, handler.BrokerList)

	// Importer
	s.registerHandler("/importer/import", POST, middleware.Protected(), handler.ImporterImport)
	s.registerHandler("/importer/apply", GET, middleware.Protected(), handler.ImporterApply)

	//	Transaction
	s.registerHandler("/transactions/get", GET, middleware.Protected(), handler.TransactionsGet)
}

func (s *server) registerHandler(path string, method HttpMethod, handlers ...fiber.Handler) {
	switch method {
	case GET:
		s.api.Get(path, handlers...)
	case POST:
		s.api.Post(path, handlers...)
	case PATCH:
		s.api.Patch(path, handlers...)
	case DELETE:
		s.api.Delete(path, handlers...)
	}
}

func (s *server) start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.port))
}
