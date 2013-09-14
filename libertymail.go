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

	// Start connection service
	log.Println("Starting connection service")
	connChan := make(chan net.Conn)
	closeChan := make(chan bool)

	go grid.Listen(uint16(*port), connChan, closeChan)

	// Start command service
	log.Println("Starting command service")
	cmdChan := make(chan string)

	go api.Console(cmdChan)

L1:
	for { // Event loop

		select {
		case connection := <-connChan:

			grid.Handshake(peers, connection, false)

		case command := <-cmdChan:

			log.Printf("Received command %s\n", command)

			if strings.HasPrefix(command, "QUIT") {

				break L1

			} else if strings.HasPrefix(command, "CONNECT") {

				items := strings.Split(command, " ")

				if len(items) > 1 {

					grid.Connect(peers, items[1])

				} else {

					log.Printf("Invalid command: %s", command)
				}
			}
		}
	}

	log.Println("Shutting down connection service")

	closeChan <- true
	<-closeChan

	for k, v := range peers {

		log.Println("Closing peer", k)
		v.Close()
	}

	log.Println("Done.")
}
