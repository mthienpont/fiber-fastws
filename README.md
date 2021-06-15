# fiber-fastws

## Installation

Lightweight wrapper for the `dgrr/fastws` library for Fiber.

```
go get -u github.com/mthienpont/fiber-fastws
```

## Example

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/dgrr/fastws"
	"github.com/mthienpont/fiber-fastws"
)

func asyncHandler(conn *websocket.Conn) {
	dataChannel := make(chan []byte)
	var b []byte

	go func() {
		for {
			_, msg, err := conn.ReadMessage(b)
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			dataChannel <- msg
		}
	}()

	for data := range dataChannel {
		conn.WriteMessage(fastws.ModeBinary, data)
	}
}

func main() {
	app := fiber.New()
	app.Get("/ws", websocket.New(asyncHandler))

	app.Listen(":3000") // ws://localhost:3000/ws
}
```



The feature comparison and benchmarks below were taken from [here](https://github.com/dgrr/fastws).

>## fastws vs gorilla vs nhooyr vs gobwas
>
>| Features | [fastws](https://github.com/dgrr/fastws) | [Gorilla](https://github.com/savsgio/websocket)| [Nhooyr](https://github.com/nhooyr/websocket) | [gowabs](https://github.>com/gobwas/ws) |
>| --- | --- | --- | --- | --- |
>| Concurrent R/W                          | Yes            | No           | No. Only writes | No           |
>| Passes Autobahn Test Suite              | Mostly         | Yes          | Yes             | Mostly       |
>| Receive fragmented message              | Yes            | Yes          | Yes             | Yes          |
>| Send close message                      | Yes            | Yes          | Yes             | Yes          |
>| Send pings and receive pongs            | Yes            | Yes          | Yes             | Yes          |
>| Get the type of a received data message | Yes            | Yes          | Yes             | Yes          |
>| Compression Extensions                  | On development | Experimental | Yes             | No (?)       |
>| Read message using io.Reader            | Not planned    | Yes          | No              | No (?)       |
>| Write message using io.WriteCloser      | Not planned    | Yes          | No              | No (?)       |
>
>## Benchmarks: fastws vs gorilla vs nhooyr vs gobwas
>
>Fastws:
>```
>$ go test -bench=Fast -benchmem -benchtime=10s
>Benchmark1000FastClientsPer10Messages-8          225367248    52.6 ns/op       0 B/op   0 allocs/op
>Benchmark1000FastClientsPer100Messages-8        1000000000     5.48 ns/op      0 B/op   0 allocs/op
>Benchmark1000FastClientsPer1000Messages-8       1000000000     0.593 ns/op     0 B/op   0 allocs/op
>Benchmark100FastMsgsPerConn-8                   1000000000     7.38 ns/op      0 B/op   0 allocs/op
>Benchmark1000FastMsgsPerConn-8                  1000000000     0.743 ns/op     0 B/op   0 allocs/op
>Benchmark10000FastMsgsPerConn-8                 1000000000     0.0895 ns/op    0 B/op   0 allocs/op
>Benchmark100000FastMsgsPerConn-8                1000000000     0.0186 ns/op    0 B/op   0 allocs/op
>```
>
>Gorilla:
>```
>$ go test -bench=Gorilla -benchmem -benchtime=10s
>Benchmark1000GorillaClientsPer10Messages-8       128621386    97.5 ns/op      86 B/op   1 allocs/op
>Benchmark1000GorillaClientsPer100Messages-8     1000000000    11.0 ns/op       8 B/op   0 allocs/op
>Benchmark1000GorillaClientsPer1000Messages-8    1000000000     1.12 ns/op      0 B/op   0 allocs/op
>Benchmark100GorillaMsgsPerConn-8                 849490059    14.0 ns/op       8 B/op   0 allocs/op
>Benchmark1000GorillaMsgsPerConn-8               1000000000     1.42 ns/op      0 B/op   0 allocs/op
>Benchmark10000GorillaMsgsPerConn-8              1000000000     0.143 ns/op     0 B/op   0 allocs/op
>Benchmark100000GorillaMsgsPerConn-8             1000000000     0.0252 ns/op    0 B/op   0 allocs/op
>```
>
>Nhooyr:
>```
>$ go test -bench=Nhooyr -benchmem -benchtime=10s
>Benchmark1000NhooyrClientsPer10Messages-8        121254158   114 ns/op        87 B/op   1 allocs/op
>Benchmark1000NhooyrClientsPer100Messages-8      1000000000    11.1 ns/op       8 B/op   0 allocs/op
>Benchmark1000NhooyrClientsPer1000Messages-8     1000000000     1.19 ns/op      0 B/op   0 allocs/op
>Benchmark100NhooyrMsgsPerConn-8                  845071632    15.1 ns/op       8 B/op   0 allocs/op
>Benchmark1000NhooyrMsgsPerConn-8                1000000000     1.47 ns/op      0 B/op   0 allocs/op
>Benchmark10000NhooyrMsgsPerConn-8               1000000000     0.157 ns/op     0 B/op   0 allocs/op
>Benchmark100000NhooyrMsgsPerConn-8              1000000000     0.0251 ns/op    0 B/op   0 allocs/op
>```
>
>Gobwas:
>```
>$ go test -bench=Gobwas -benchmem -benchtime=10s
>Benchmark1000GobwasClientsPer10Messages-8         98497042   106 ns/op        86 B/op   1 allocs/op
>Benchmark1000GobwasClientsPer100Messages-8      1000000000    13.4 ns/op       8 B/op   0 allocs/op
>Benchmark1000GobwasClientsPer1000Messages-8     1000000000     1.19 ns/op      0 B/op   0 allocs/op
>Benchmark100GobwasMsgsPerConn-8                  833576667    14.6 ns/op       8 B/op   0 allocs/op
>Benchmark1000GobwasMsgsPerConn-8                1000000000     1.46 ns/op      0 B/op   0 allocs/op
>Benchmark10000GobwasMsgsPerConn-8               1000000000     0.156 ns/op     0 B/op   0 allocs/op
>Benchmark100000GobwasMsgsPerConn-8              1000000000     0.0262 ns/op    0 B/op   0 allocs/op
>```
>
>The source files are in [this](https://github.com/dgrr/fastws/tree/master/stress-tests/) folder.

