package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func StartServer() {
	socket, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}

		log.Println("server: accepted connection")

		go func(conn net.Conn) {
			buf := make([]byte, 1024)

			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Fatalf("Error reading data: %v", err)
				}

				log.Printf("server: read %q", buf[:n])

				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Fatalf("error writing data: %v", err)
				}
			}
		}(conn)
	}
}

func main() {

	go StartServer()

	// simulate server startup time
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	for i := 0; i < 3; i++ {
		_, err := conn.Write([]byte(fmt.Sprintf("hello %d", i)))
		if err != nil {
			fmt.Println("Error sending message:", err)
			os.Exit(1)
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		log.Printf("client: read %q", buf[:n])
	}
}
