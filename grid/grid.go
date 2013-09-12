// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package grid

import (
	"fmt"
	"log"
	"net"
	"time"
)

func Listen(port uint16, connChan chan<- net.Conn, closeChan chan bool) {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if nil != err {

		log.Fatalln("grid.Listen:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {

		log.Fatalln("grid.Listen:", err)
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

				log.Fatalln("grid.Listen:", err)
				return
			}
		}

		// Send connection back to client
		connChan <- conn
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
