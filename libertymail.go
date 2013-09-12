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

func Listen(port uint16, connChan chan<- net.Conn, closeChan chan bool) {

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

func Command(cmdChan chan<- string) {

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

		// If we have a quit command, exit this fiber
		if strings.HasPrefix(cmd, "QUIT") {
			break
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
	w := bufio.NewWriter(fd)
	log.SetOutput(w)
	defer fd.Close()
	defer w.Flush()

	peers := make(map[string]net.Conn)

	// Start connection listener
	log.Println("Starting connection listener")
	connChan := make(chan net.Conn)
	closeChan := make(chan bool)
	go Listen(uint16(*port), connChan, closeChan)

	// Start command listener
	log.Println("Starting command listener")
	cmdChan := make(chan string)
	go Command(cmdChan)

L1:
	for { // Event loop
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

	log.Println("Shutting down connection listener...")
	closeChan <- true
	<-closeChan

	for k, v := range peers {
		log.Println("Closing peer", k)
		v.Close()
	}

	log.Println("Done.")
}
