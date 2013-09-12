// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package api

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Console(cmdChan chan<- string) {

	// Read commands from stdin
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("LM: ")

		line, err := reader.ReadString('\n')

		if err != nil {

			log.Fatalln("CommandService:", err)
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
