// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package grid

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ListenService struct {
	Port           uint16
	ConnectionChan chan net.Conn
	CloseChan      chan struct{}
	ServiceGroup   *sync.WaitGroup
}

func (ls *ListenService) Run() {

	log.Println("Starting listen service")

	ls.ServiceGroup.Add(1)
	defer ls.ServiceGroup.Done()

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("127.0.0.1:%d", ls.Port))
	if nil != err {

		log.Fatalln("grid.ListenService.Run:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {

		log.Fatalln("grid.ListenService.Run:", err)
		return
	}

	log.Println("Listening on", listener.Addr())

	for {
		// Shutdown gracefully if we have a close signal waiting
		select {
		case <-ls.CloseChan:

			listener.Close()
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

				log.Fatalln("grid.ListenService.Run:", err)
				return
			}
		}

		// Send connection back to client
		ls.ConnectionChan <- conn
	}
}

type ConnectService struct {
	AddressChan    chan string
	ConnectionChan chan net.Conn
	CloseChan      chan struct{}
	ServiceGroup   *sync.WaitGroup
}

func (cs *ConnectService) Run() {

	log.Println("Starting connect service")

	cs.ServiceGroup.Add(1)
	defer cs.ServiceGroup.Done()

	for {
		select {

		case <-cs.CloseChan:

			return

		case addr := <-cs.AddressChan:

			conn, err := net.Dial("tcp", addr)

			if err != nil {

				log.Println("grid.ConnectService.Run:", err)

			} else {

				cs.ConnectionChan <- conn
			}
		}
	}
}

type HandshakeService struct {
	ConnectionChan chan net.Conn
	CloseChan      chan struct{}
	ServiceGroup   *sync.WaitGroup
}

func (hs *HandshakeService) Run() {

	log.Println("Starting handshake service")

	hs.ServiceGroup.Add(1)
	defer hs.ServiceGroup.Done()

	for {
		select {

		case <-hs.CloseChan:

			return

		case connection := <-hs.ConnectionChan:

			log.Println("Doing handshake with", connection.RemoteAddr())
			// TODO: Handshaking
			hs.ConnectionChan <- connection
		}
	}
}

type InitiateHandshakeService struct {
	ConnectionChan chan net.Conn
	CloseChan      chan struct{}
	ServiceGroup   *sync.WaitGroup
}

func (ihs *InitiateHandshakeService) Run() {

	log.Println("Starting initiate handshake service")

	ihs.ServiceGroup.Add(1)
	defer ihs.ServiceGroup.Done()

	for {
		select {

		case <-ihs.CloseChan:

			return

		case connection := <-ihs.ConnectionChan:

			log.Println("Initiating handshake with", connection.RemoteAddr())
			// TODO: Handshaking
			ihs.ConnectionChan <- connection
		}
	}
}
