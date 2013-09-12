// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	//"libertymail-go/proto"
)

func ListenService(port uint16, connChan chan<- net.Conn, closeChan chan bool) {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if nil != err {
		log.Fatalln("Listen:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln("Listen:", err)
		return
	}

	log.Println("Listening on", listener.Addr())

	for {
		// Shutdown gracefully if we have a close signal waiting
		select {
		case <-closeChan:
			log.Println("Stopping listening on", listener.Addr())
			listener.Close()
			closeChan <- true
			return
		default:
		}

		// Accept incoming connections for 3 seconds
		listener.SetDeadline(time.Now().Add(time.Second * 3))

		conn, err := listener.AcceptTCP()
		if err != nil {
			// Accept returns an error on timeout. If we have a timeout, continue
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			} else {
				log.Fatalln("Listen:", err)
				return
			}
		}

		// Send connection back to client
		connChan <- conn
	}
}

func CommandService(cmdChan chan<- string) {

	// Read commands from stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("LM: ")

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Command:", err)
			break
		}

		cmd := strings.ToUpper(strings.Trim(line, "\n\r\t "))

		// Send command back to client
		cmdChan <- cmd

		// If we have a quit command, exit this service
		if strings.HasPrefix(cmd, "QUIT") {
			return
		}
	}
}

func Handshake(peers map[string]net.Conn, conn net.Conn, initiate bool) {

	log.Println("Handshaking with", conn.RemoteAddr())
	// FIXME: do handshaking
	peers[conn.RemoteAddr().String()] = conn
}

func Connect(peers map[string]net.Conn, addr string) bool {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Connect:", err)
		return false
	} else {
		log.Println("Connecting to", conn.RemoteAddr())
		Handshake(peers, conn, true)
	}

	return true
}

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

	w := bufio.NewWriter(fd)
	log.SetOutput(w)
	defer w.Flush()

	// Temporary map to hold peer connections
	peers := make(map[string]net.Conn)

	// Start connection service
	log.Println("Starting connection service")
	connChan := make(chan net.Conn)
	closeChan := make(chan bool)
	go ListenService(uint16(*port), connChan, closeChan)

	// Start command service
	log.Println("Starting command service")
	cmdChan := make(chan string)
	go CommandService(cmdChan)

L1:
	for { // Event loop

		w.Flush() // Give those log files a little push

		select {

		case connection := <-connChan:

			Handshake(peers, connection, false)
		case command := <-cmdChan:

			log.Printf("Received command %s\n", command)

			if strings.HasPrefix(command, "QUIT") {
				break L1
			} else if strings.HasPrefix(command, "CONNECT") {

				items := strings.Split(command, " ")
				if len(items) > 1 {
					Connect(peers, items[1])
				} else {
					log.Printf("Invalid command: %s", command)
				}
			}
		}
	}

	log.Println("Shutting down connection service...")

	closeChan <- true
	<-closeChan

	for k, v := range peers {
		log.Println("Closing peer", k)
		v.Close()
	}

	log.Println("Done.")
}
