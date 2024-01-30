package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	conf "github.com/geoffjay/plantd/app/config"
	"github.com/geoffjay/plantd/app/handlers"
	"github.com/geoffjay/plantd/core/util"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	log "github.com/sirupsen/logrus"
)

type service struct{}

func (s *service) init() {
	log.WithFields(log.Fields{
		"service": "app",
		"context": "service.init",
	}).Debug("initializing")
}

func (s *service) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	log.WithFields(log.Fields{
		"service": "app",
		"context": "service.run",
	}).Debug("starting")

	wg.Add(1)
	go s.runApp(ctx, wg)

	<-ctx.Done()

	log.WithFields(log.Fields{
		"service": "app",
		"context": "service.run",
	}).Debug("exiting")
}

func (s *service) runApp(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	config := conf.GetConfig()

	fields := log.Fields{"service": "app", "context": "service.run-app"}
	bindAddress := util.Getenv("PLANTD_APP_BIND_ADDRESS", "0.0.0.0")
	bindPort, err := strconv.Atoi(util.Getenv("PLANTD_APP_BIND_PORT", "8080"))
	if err != nil {
		log.WithFields(fields).Fatal(err)
	}

	log.WithFields(fields).Debug("starting server")

	go func() {
		engine := html.New("app/views", ".tmpl")
		engine.Reload(true)
		if config.Env == "development" {
			engine.Debug(true)
		}

		app := fiber.New(fiber.Config{
			Views:       engine,
			JSONEncoder: json.Marshal,
			JSONDecoder: json.Unmarshal,
		})

		handlers.SessionStore = session.New(session.Config{
			CookieSecure: true,
		})

		// csrfConfig := csrf.Config{
		// 	KeyLookup:         "header:" + csrf.HeaderName,
		// 	CookieName:        "__Host-csrf_",
		// 	CookieSameSite:    "Lax",
		// 	CookieSecure:      true,
		// 	CookieSessionOnly: true,
		// 	CookieHTTPOnly:    true,
		// 	Expiration:        1 * time.Hour,
		// 	KeyGenerator:      futils.UUIDv4,
		// 	// ErrorHandler:      DefaultErrorHandler,
		// 	Extractor:         csrf.CsrfFromHeader(csrf.HeaderName),
		// 	Session:           store,
		// 	SessionKey:        "fiber.csrf.token",
		// 	HandlerContextKey: "fiber.csrf.handler",
		// }

		corsConfig := cors.Config{
			AllowCredentials: true,
			AllowOrigins:     "*",
			AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		}

		app.Use(helmet.New())
		// app.Use(csrf.New(csrfConfig))
		app.Use(cors.New(corsConfig))
		app.Use(logger.New())
		app.Use(recover.New())
		app.Use(etag.New())
		app.Use(limiter.New(limiter.Config{
			Expiration: 30 * time.Second,
			Max:        50,
		}))

		initializeRouter(app)

		log.Fatal(app.Listen(fmt.Sprintf("%s:%d", bindAddress, bindPort)))
	}()

	<-ctx.Done()

	log.WithFields(fields).Debug("exiting server")
}
