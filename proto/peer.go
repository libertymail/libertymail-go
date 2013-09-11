package proto

import (
//"bytes"
//"errors"
)

type peer struct {
	IP        string
	Port      uint16
	PublicKey []byte
	Version   byte
}

func NewPeer() *peer {

	return new(peer)
}

func NewPeerFrom(ip string, port uint16) *peer {

	peer := new(peer)
	peer.IP = ip
	peer.Port = port

	return peer
}

func (p *peer) Serialize() ([]byte, error) {
}

func (p *peer) Deserialize(packet []byte) error {
}
