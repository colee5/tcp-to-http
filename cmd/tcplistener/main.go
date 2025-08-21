package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const BYTE_SIZE = 8

func main() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	listener, error := net.Listen("tcp", ":42069")
	if error != nil {
		log.Fatal(error)
	}

	defer listener.Close()
	for {
		// Wait for a connection.
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Accepted connection")
		output := getLinesChannel(conn)

		for line := range output {
			fmt.Printf("read: %s\n", line)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			log.Println("Closed connection")
			c.Close()
		}(conn)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	buffer := make([]byte, BYTE_SIZE)

	currentLine := "" // Accumulates text until we hit a complete line

	go func() {
		// This runs in the BACKGROUND
		// while the function returns immediately
		defer close(lines)

		for {
			n, err := f.Read(buffer)

			if n > 0 {
				chunk := string(buffer[:n])         // Convert 8 bytes to string
				parts := strings.Split(chunk, "\n") // Split on newlines

				for i := 0; i < len(parts)-1; i++ {
					// Complete the current line and send it to the channel
					completeLine := currentLine + parts[i]
					lines <- completeLine
					// Reset current line for the next line
					currentLine = ""
				}

				currentLine += parts[len(parts)-1]

			}
			if err != nil {
				if err == io.EOF {
					// Send any remaining data as the last line
					if currentLine != "" {
						lines <- currentLine
					}
					break // End of file
				}
				log.Fatal(err)
			}
		}
	}()

	return lines
}
