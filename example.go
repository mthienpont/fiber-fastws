package websocket

import (
	"fmt"
	"log"

	"github.com/dgrr/fastws"
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
	okChannel := make(chan bool)

	go listenToIncoming(dataChannel, conn)
	go echoReply(dataChannel, conn)

	working := <-okChannel
	log.Println(working)
}

func listenToIncoming(dataChannel chan []byte, conn *Conn) {
	for {
		_, msg, err := conn.ReadMessage(nil)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", msg)
		dataChannel <- msg
	}
}

func echoReply(dataChannel chan []byte, conn *Conn) {
	for data := range dataChannel {
		conn.WriteMessage(fastws.ModeBinary, data)
	}
}

func main() {
	app := fiber.New()
	app.Get("/ws", New(asyncHandler))

	app.Listen(":3000") // ws://localhost:3000/ws
}
