package websocket

import (
	"errors"
	"sync"
	"time"

	"github.com/dgrr/fastws"
	"github.com/fasthttp/websocket"
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
func New(handler func(*Conn), config ...Config) fiber.Handler {
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

	return func(c *fiber.Ctx) error {
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

		return nil
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

// Constants are taken from https://github.com/fasthttp/websocket/blob/master/conn.go#L43

// Close codes defined in RFC 6455, section 11.7.
const (
	CloseNormalClosure           = 1000
	CloseGoingAway               = 1001
	CloseProtocolError           = 1002
	CloseUnsupportedData         = 1003
	CloseNoStatusReceived        = 1005
	CloseAbnormalClosure         = 1006
	CloseInvalidFramePayloadData = 1007
	ClosePolicyViolation         = 1008
	CloseMessageTooBig           = 1009
	CloseMandatoryExtension      = 1010
	CloseInternalServerErr       = 1011
	CloseServiceRestart          = 1012
	CloseTryAgainLater           = 1013
	CloseTLSHandshake            = 1015
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

var (
	ErrBadHandshake = errors.New("websocket: bad handshake")
	ErrCloseSent    = errors.New("websocket: close sent")
	ErrReadLimit    = errors.New("websocket: read limit exceeded")
)

// IsWebSocketUpgrade returns true if the client requested upgrade to the
// WebSocket protocol.
func IsWebSocketUpgrade(c *fiber.Ctx) bool {
	return websocket.FastHTTPIsWebSocketUpgrade(c.Context())
}
