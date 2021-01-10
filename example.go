package websocket

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func wsHandler(conn *Conn) {
	fmt.Fprintf(conn, "Hello world")
	var b []byte
	for {
		mt, msg, err := conn.ReadMessage(b)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		_, err = conn.WriteMessage(mt, msg)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	app := fiber.New()
	app.Get("/ws", New(wsHandler))

	app.Listen(":3000") // ws://localhost:3000/ws
}
