package proto

import (
//"bytes"
//"errors"
)

type peer struct {
	IP        []byte
	Port      uint16
	PublicKey []byte
	Version   byte
}

func NewPeer() *peer {

	return new(peer)
}

func (p *peer) Serialize() ([]byte, error) {
}

func (p *peer) Deserialize(packet []byte) error {
}
