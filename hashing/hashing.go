// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package hashing

import (	
	"crypto/sha512"

	"libertymail-go/ripemd160"
)

func SHA512(p []byte) []byte {

	sha := sha512.New()
	sha.Write(p)	
	return sha.Sum(nil)
}

func SHA512x2(p []byte) []byte {

	sha1, sha2 := sha512.New(), sha512.New()
	sha1.Write(p)
	sha2.Write(sha1.Sum(nil))
	return sha2.Sum(nil)
}

func RIPEMD160(p []byte) []byte {

	ripe := ripemd160.New()
	ripe.Write(p)	
	return ripe.Sum(nil)
}

func RIPEMD160x2(p []byte) []byte {

	ripe1, ripe2 := ripemd160.New(), ripemd160.New()
	ripe1.Write(p)
	ripe2.Write(ripe1.Sum(nil))
	return ripe2.Sum(nil)
}