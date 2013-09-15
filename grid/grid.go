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

type ListenService struct {
	Port           uint16
	ConnectionChan chan net.Conn
	CloseChan      chan bool
}

func (ls *ListenService) Run() {

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", ls.Port))
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
		case <-ls.CloseChan:

			log.Println("Stopping listen service on", listener.Addr())
			listener.Close()
			ls.CloseChan <- true
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

				log.Fatalln("grid.ListenService:", err)
				return
			}
		}

		// Send connection back to client
		ls.ConnectionChan <- conn
	}
}

func (ls *ListenService) Close() {

	ls.CloseChan <- true
	<-ls.CloseChan
}

type ConnectService struct {
	AddressChan    chan string
	ConnectionChan chan net.Conn
	CloseChan      chan bool
}

func (cs *ConnectService) Run() {

	for {
		select {

		case <-cs.CloseChan:

			log.Println("Stopping connection service")
			cs.CloseChan <- true
			return

		case addr := <-cs.AddressChan:

			conn, err := net.Dial("tcp", addr)

			if err != nil {
				log.Println("grid.ConnectService:", err)
			} else {
				cs.ConnectionChan <- conn
			}
		}
	}
}

func (cs *ConnectService) Close() {

	cs.CloseChan <- true
	<-cs.CloseChan
}

type HandshakeService struct {
	ConnectionChan chan net.Conn
	CloseChan      chan bool
}

func (hs *HandshakeService) Run() {

	for {
		select {

		case <-hs.CloseChan:

			log.Println("Stopping handshake service")
			hs.CloseChan <- true
			return

		case connection := <-hs.ConnectionChan:

			log.Println("Doing handshake with", connection.RemoteAddr())
			// TODO: Handshaking
			hs.ConnectionChan <- connection
		}
	}
}

func (hs *HandshakeService) Close() {

	hs.CloseChan <- true
	<-hs.CloseChan
}

type InitiateHandshakeService struct {
	ConnectionChan chan net.Conn
	CloseChan      chan bool
}

func (ihs *InitiateHandshakeService) Run() {

	for {
		select {

		case <-ihs.CloseChan:

			log.Println("Stopping initiate handshake service")
			ihs.CloseChan <- true
			return

		case connection := <-ihs.ConnectionChan:

			log.Println("Initiating handshake with", connection.RemoteAddr())
			// TODO: Handshaking
			ihs.ConnectionChan <- connection
		}
	}
}

func (ihs *InitiateHandshakeService) Close() {

	ihs.CloseChan <- true
	<-ihs.CloseChan
}
