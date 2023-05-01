package main

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
    "fmt"

	"github.com/Andreaswiv/funtemps/conv"
	"github.com/Andreaswiv/is105sem03/mycrypt"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.3:8888")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}

			go handleConnection(conn)
		}
	}()

	wg.Wait()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return
	}

	msg := string(buf[:n])
	switch {
	case strings.HasPrefix(msg, "ping"):
		_, err = conn.Write([]byte("pong"))
	case strings.HasPrefix(msg, "kjevik"):
		parts := strings.Split(msg, ";")
		if len(parts) != 4 {
			log.Println("invalid message format")
			return
		}

		degrees, err := strconv.Atoi(parts[3])
		if err != nil {
			log.Println(err)
			return
		}

		fahrenheit := conv.CelsiusToFahrenheit(float64(degrees))
		response := fmt.Sprintf("%s;%s;%s;%.1f", parts[0], parts[1], parts[2], fahrenheit)
		_, err = conn.Write([]byte(response))
	default:
		dekryptertMelding := mycrypt.Krypter([]rune(msg), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
		log.Println("Dekryptert melding: ", string(dekryptertMelding))
		_, err = conn.Write([]byte(string(dekryptertMelding)))
	}

	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return
	}
}
