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

func Listen(port uint16, c chan<- net.Conn, shutdown chan bool) {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if nil != err {
		log.Fatalln(err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())

	for {

		select {
		case <-shutdown:
			log.Println("stopping listening on", listener.Addr())
			listener.Close()
			shutdown <- true
			return
		default:
		}

		listener.SetDeadline(time.Now().Add(time.Second * 3))
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		c <- conn
	}
}

func Command(c chan<- string) {

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("LM: ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("CommandManager: " + err.Error())
			break
		}

		c <- strings.ToUpper(strings.Trim(line, "\n\r\t "))
	}
}

func Handshake(peers map[string]net.Conn, conn net.Conn) {

	log.Println("Handshaking with", conn.RemoteAddr().String())
	peers[conn.RemoteAddr().String()] = conn
}

func Connect(peers map[string]net.Conn, addr string) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Connect error: " + err.Error())
	} else {
		log.Println("Connecting to", conn.RemoteAddr().String())
		peers[conn.RemoteAddr().String()] = conn
	}
}

func main() {

	port := flag.Uint("port", 30000, "The listening port")
	logfile := flag.String("logfile", "log.txt", "The log file")
	flag.Parse()

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
	shutdown := make(chan bool)
	go Listen(uint16(*port), connChan, shutdown)

	// Start command listener
	log.Println("Starting command listener")
	cmdChan := make(chan string)
	go Command(cmdChan)

L1:
	for {
		select {

		case connection := <-connChan:
			Handshake(peers, connection)

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

	log.Println("Shutting down connection manager...")
	shutdown <- true
	<-shutdown

	for k, v := range peers {
		log.Println("Closing peer", k)
		v.Close()
	}

	log.Println("Done.")
}
