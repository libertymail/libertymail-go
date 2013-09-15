// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"libertymail-go/api"
	"libertymail-go/grid"
)

func main() {

	// Parse commandline flags
	port := flag.Uint("port", 30000, "The listening port")
	logfile := flag.String("logfile", "log.txt", "The log file")
	flag.Parse()

	if flag.NFlag() < 2 {
		fmt.Println("Usage:", os.Args[0], "-port=... -logfile=...")
		return
	}

	// Setup logfile
	fd, err := os.Create(*logfile)
	if err != nil {

		panic(err)
	}
	defer fd.Close()
	log.SetOutput(fd)

	// Temporary map to hold peer connections
	peers := make(map[string]net.Conn)

	// Start command service
	log.Println("Starting command service")
	consoleService := &api.ConsoleService{make(chan string)}
	go consoleService.Run()

	// Start listen service
	log.Println("Starting listen service")
	listenService := &grid.ListenService{uint16(*port), make(chan net.Conn), make(chan bool)}
	go listenService.Run()

	// Start connect service
	log.Println("Starting connect service")
	connectService := &grid.ConnectService{make(chan string), make(chan net.Conn), make(chan bool)}
	go connectService.Run()

	// Start handshake service
	log.Println("Starting handshake service")
	handshakeService := &grid.HandshakeService{make(chan net.Conn), make(chan bool)}
	go handshakeService.Run()

	// Start initiate handshake service
	log.Println("Starting initiate handshake service")
	initiateHandshakeService := &grid.InitiateHandshakeService{make(chan net.Conn), make(chan bool)}
	go initiateHandshakeService.Run()

L1:
	for { // Event loop

		select {
		case command := <-consoleService.CommandChan:

			log.Printf("Received command %s\n", command)

			if strings.HasPrefix(command, "QUIT") {

				break L1

			} else if strings.HasPrefix(command, "CONNECT") {

				items := strings.Split(command, " ")
				if len(items) > 1 {

					connectService.AddressChan <- items[1]

				} else {

					log.Printf("Invalid command: %s", command)
				}
			} else if strings.HasPrefix(command, "LIST") {

				log.Println("Connected peers:")
				for k, _ := range peers {

					log.Println(k)
				}
			}

		case connection := <-listenService.ConnectionChan:

			handshakeService.ConnectionChan <- connection

		case connection := <-handshakeService.ConnectionChan:

			peers[connection.RemoteAddr().String()] = connection
			//db.RegisterPeer(connection)

		case connection := <-connectService.ConnectionChan:

			initiateHandshakeService.ConnectionChan <- connection

		case connection := <-initiateHandshakeService.ConnectionChan:

			peers[connection.RemoteAddr().String()] = connection
			//db.RegisterPeer(connection)
		}
	}

	listenService.Close()
	connectService.Close()
	handshakeService.Close()
	initiateHandshakeService.Close()

	for k, v := range peers {

		log.Println("Closing peer", k)
		v.Close()
	}
}
