package websocket

import (
	"time"

	"github.com/dgrr/fastws"
	"github.com/gofiber/fiber/v2"
)

// Config ...
type Config struct {
	Filter                          func(*fiber.Ctx) bool
	HandshakeTimeout                time.Duration
	Subprotocols                    []string
	Origins                         []string
	ReadBufferSize, WriteBufferSize int
	EnableCompression               bool
}

// New takes a fastws handler and upgrades the connection
func New(handler func(conn *fastws.Conn), config ...Config) fiber.Handler {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if len(cfg.Origins) == 0 {
		cfg.Origins = []string{"*"}
	}

	var upgrader = fastws.Upgrader{
		Handler:   handler,
		Protocols: cfg.Subprotocols,
		Compress:  cfg.EnableCompression,
	}

	return func(c *fiber.Ctx) error {
		upgrader.Upgrade(c.Context())
		return nil
	}
}
