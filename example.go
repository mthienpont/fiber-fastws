package websocket

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func wsHandler(conn *Conn) {
	fmt.Fprintf(conn, "Hello world")

	for {
		mt, msg, err := conn.ReadMessage(nil)
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

func asyncHandler(conn *Conn) {
	dataChannel := make(chan []byte)

	go func() {
		for {
			_, msg, err := conn.ReadMessage(nil)
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			dataChannel <- msg
		}
	}()

	for data := range dataChannel {
		conn.WriteMessage(1, data)
	}
}

func main() {
	app := fiber.New()
	app.Get("/ws", New(asyncHandler))

	app.Listen(":3000") // ws://localhost:3000/ws
}
