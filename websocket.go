package websocket

import (
	"sync"
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
func New(handler func(*Conn), config ...Config) func(c *fiber.Ctx) {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if len(cfg.Origins) == 0 {
		cfg.Origins = []string{"*"}
	}

	// var upgrader = fastws.Upgrader{
	// 	Handler:   handler,
	// 	Protocols: cfg.Subprotocols,
	// 	Compress:  cfg.EnableCompression,
	// }

	return func(c *fiber.Ctx) {
		conn := acquireConn()
		// locals
		c.Context().VisitUserValues(func(key []byte, value interface{}) {
			conn.locals[string(key)] = value
		})
		// queries
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			conn.queries[string(key)] = string(value)
		})
		// upgrade
		upgrade := fastws.Upgrade(func(fconn *fastws.Conn) {
			conn.Conn = fconn
			defer releaseConn(conn)
			handler(conn)
		})

		upgrade(c.Context())
	}
}

// Conn ...
type Conn struct {
	*fastws.Conn
	locals  map[string]interface{}
	params  map[string]string
	cookies map[string]string
	queries map[string]string
}

// Conn pool
var poolConn = sync.Pool{
	New: func() interface{} {
		return new(Conn)
	},
}

// Acquire Conn from pool
func acquireConn() *Conn {
	conn := poolConn.Get().(*Conn)
	conn.locals = make(map[string]interface{})
	conn.params = make(map[string]string)
	conn.queries = make(map[string]string)
	conn.cookies = make(map[string]string)
	return conn
}

// Return Conn to pool
func releaseConn(conn *Conn) {
	conn.Conn = nil
	poolConn.Put(conn)
}

// Locals ...
func (conn *Conn) Locals(key string) interface{} {
	return conn.locals[key]
}

// Query ...
func (conn *Conn) Query(key string, defaultValue ...string) string {
	v, ok := conn.queries[key]
	if !ok && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return v
}
