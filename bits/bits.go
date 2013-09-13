// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package bits

import (
	"crypto/sha256"
	"crypto/rand"
	"errors"
	"io"

	"libertymail-go/ripemd160"
)

func Checksum(p []byte, n int) ([]byte, error) {

	if n > len(p) {
		return nil, errors.New("bits.Checksum: checksum size is bigger than digest length")
	}

	npad := n - (len(p) % n)
	check := make([]byte, n)
	buf := make([]byte, len(p)+npad)
	copy(buf, p)

	for i := 0; i < len(buf); i += n {
		for j := 0; j < n; j++ {
			check[j] ^= buf[i+j]
		}
	}

	return check, nil
}

/* Generates a block of 192 random bits */
func GenerateRandomBlock() ([]byte, error) {
	b := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, b)
	if len(b) != 24 || err != nil {
		b = nil
	}
	return b, err
}

func SHA256(p []byte) []byte {

	sha := sha256.New()
	sha.Write(p)

	return sha.Sum(nil)
}

func RIPEMD160(p []byte) []byte {

	ripe := ripemd160.New()
	ripe.Write(p)

	return ripe.Sum(nil)
}

func Hash(p []byte) []byte {

	return RIPEMD160(SHA256(p))
}
