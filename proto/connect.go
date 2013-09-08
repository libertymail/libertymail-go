// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package proto

import (
//"bytes"
//"errors"
)

type connectMessage struct {
	Timestamp int64
}

func NewConnectMessage() *connectMessage {

	return new(connectMessage)
}

func (cm *connectMessage) Serialize() ([]byte, error) {
}

func (cm *connectMessage) Deserialize(packet []byte) error {
}
