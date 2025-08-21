package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resolved address: %s\n", addr)
	fmt.Printf("IP: %s, Port: %d\n", addr.IP, addr.Port)

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">\n")

		fmt.Print("Enter message: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(userInput))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Message sent!")
	}
}
