// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package main

import (
	"fmt"

	"libertymail-go/address"
)

func main() {

	addr, _ := address.NewAddress(1, 0, false)

	fmt.Println(addr.Identifier)
}
