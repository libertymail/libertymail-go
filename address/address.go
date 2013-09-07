// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package address

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"strings"

	"libertymail-go/base58"
	"libertymail-go/hashing"
)

type address struct {
	Version, Privacy byte
	Identifier       string
	Key              *ecdsa.PrivateKey
}

func NewAddress(version, privacy byte) (*address, error) {

	var err error
	addr := new(address)
	addr.Key, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, errors.New("address.NewAddress: Error generating ecdsa keys: " + err.Error())
	}

	addr.Version = version
	addr.Privacy = privacy

	var ident bytes.Buffer
	ident.WriteByte(version)
	ident.WriteByte(privacy)

	var buf bytes.Buffer
	buf.Write(addr.Key.PublicKey.X.Bytes())
	buf.Write(addr.Key.PublicKey.Y.Bytes())
	shaDigest := hashing.SHA512x2(buf.Bytes())
	ripeDigest := hashing.RIPEMD160x2(shaDigest)

	ident.Write(ripeDigest)
	checksum := hashing.SHA512x2(ident.Bytes())[:2]
	ident.Write(checksum)

	ident58, _ := base58.Encode(ident.Bytes())
	addr.Identifier = "LM:" + ident58
	return addr, nil
}

func ValidateIdentifier(identifier string) bool {

	if len(identifier) < 27 { // LM: + version + privacy + ripe + checksum
		return false
	}

	if !strings.HasPrefix(identifier, "LM:") {
		return false
	}

	return true
}

func ValidateChecksum(identifier string) bool {

	if !ValidateIdentifier(identifier) {
		return false
	}

	raw, _ := base58.Decode(identifier[3:])
	ident := raw[:len(raw)-2]
	cs1 := raw[len(raw)-2:]
	cs2 := hashing.SHA512x2(ident)[:2]

	return bytes.Compare(cs1, cs2) == 0
}
