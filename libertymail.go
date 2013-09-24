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
	"sync"

	"libertymail-go/api"
	"libertymail-go/grid"
	"libertymail-go/proto"
	"libertymail-go/store"
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
	logfd, err := os.Create(*logfile)
	if err != nil {
		panic(err)
	}
	defer logfd.Close()
	log.SetOutput(logfd)

	// Temporary map to hold peer connections
	peers := make(map[string]net.Conn)

	// WaitGroup to synchronize service shutdown
	serviceGroup := &sync.WaitGroup{}

	// Open databases
	db := store.NewStore("./private.db", "./public.db")
	if err = db.Open(); err != nil {
		panic(err)
	}
	defer db.Close()

	addr1, err := proto.NewAddress(1, 0)
	if err != nil {
		panic(err)
	}
	db.SaveAddress(addr1)

	// Start command service
	consoleService := &api.ConsoleService{make(chan string)}
	go consoleService.Run(serviceGroup)

	// Start listen service
	listenService := &grid.ListenService{uint16(*port), make(chan net.Conn), make(chan struct{})}
	go listenService.Run(serviceGroup)

	// Start connect service
	connectService := &grid.ConnectService{make(chan string), make(chan net.Conn)}
	go connectService.Run(serviceGroup)

	// Start handshake service
	handshakeService := &grid.HandshakeService{make(chan net.Conn)}
	go handshakeService.Run(serviceGroup)

	// Start initiate handshake service
	initiateHandshakeService := &grid.InitiateHandshakeService{make(chan net.Conn)}
	go initiateHandshakeService.Run(serviceGroup)

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

	log.Println("Stopping services")

	listenService.Close()
	connectService.Close()
	handshakeService.Close()
	initiateHandshakeService.Close()

	serviceGroup.Wait()

	for k, v := range peers {

		log.Println("Closing peer", k)
		v.Close()
	}
}
