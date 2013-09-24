// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
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

	// Start console service
	apiService := &api.JsonService{make(chan string)}
	go apiService.Run(serviceGroup)

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

	var cmd api.Command
L1:
	for { // Event loop

		select {
		case request := <-apiService.StreamChan:

			if err := json.Unmarshal([]byte(request), &cmd); err != nil {
				log.Println(err)
				continue
			}

			log.Println("Received command", cmd.Name)

			switch cmd.Name {
			case "quit":

				break L1

			case "connect":

				connectService.AddressChan <- cmd.Args[0]

			case "list":

				switch cmd.Args[0] {

				case "peers":

					reply := ""
					for k, _ := range peers {
						reply += k + "\n"
					}
					apiService.StreamChan <- reply

				case "addresses":

					apiService.StreamChan <- "LM:blahblahblah\nLM:blahblahblah\n"

				default:
				}

			default:
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
