package middleware

import (
	"component-master/config"
	"component-master/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Source string

const (
	HeaderName            = "X-Csrf-Token"
	SourceCookie   Source = "cookie"
	SourceHeader   Source = "header"
	SourceURLQuery Source = "query"
)

func CorsFilter() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           86400,
	})
}

func RateLimitFiber() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               20,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}

func BasicAuthFilter(ba *config.MiddlewareConfig) fiber.Handler {
	return basicauth.New(basicauth.Config{
		Next:  nil,
		Users: map[string]string{ba.BasicAuth.Username: ba.BasicAuth.Password},
		Realm: "Forbidden",
		Authorizer: func(username, password string) bool {
			return username == ba.BasicAuth.Username && password == ba.BasicAuth.Password
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return c.SendFile(ba.Static.Unauthorized)
		},
		ContextUsername: "username",
		ContextPassword: "password",
	})
}

func CSRFFilter() fiber.Handler {
	return csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		Expiration:     30 * time.Minute,
		KeyGenerator:   util.UUIDFunc(),
	})
}
