// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

/* Json API examples

{"Name":"connect","Args":["127.0.0.1:30000"]}
{"Name":"quit","Args":null}
{"Name":"list","Args":["addresses","peers"]}

*/

type Command struct {
	Name string
	Args []string
}

func NewCommand(name string) *Command {

	c := new(Command)
	c.Name = name

	return c
}

func (c *Command) Arg(p string) *Command {

	c.Args = append(c.Args, p)

	return c
}

func (c *Command) String() string {

	b, _ := json.Marshal(c)
	return string(b)
}

type JsonService struct {
	StreamChan chan string
}

func (js *JsonService) Run(serviceGroup *sync.WaitGroup) {

	log.Println("Starting Json service")

	serviceGroup.Add(1)
	defer serviceGroup.Done()

	cs := &consoleService{make(chan string)}
	go cs.Run()

	for {
		select {

		case cmd := <-cs.CommandChan:

			js.StreamChan <- cmd

			// If we have a quit command, exit this service
			if strings.HasPrefix(cmd, "{\"Name\":\"quit\"") {

				log.Println("Quitting API service")
				return
			}

		case reply := <-js.StreamChan:

			fmt.Print(reply)
		}
	}
}

type consoleService struct {
	CommandChan chan string
}

func (cs *consoleService) Run() {

	log.Println("Starting console service")

	// Read commands from stdin
	reader := bufio.NewReader(os.Stdin)

	for {

		line, err := reader.ReadString('\n')

		if err != nil {

			fmt.Println("ERROR:", err)
			break
		}

		cmd := strings.Trim(line, "\n\r\t ")

		cs.CommandChan <- cmd

		// If we have a quit command, exit this service
		if strings.HasPrefix(cmd, "{\"Name\":\"quit\"") {

			log.Println("Quitting console service")
			return
		}
	}
}
