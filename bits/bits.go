// LICENSE: GNU General Public License version 2
// CONTRIBUTORS AND COPYRIGHT HOLDERS (c) 2013:
// Dag Robøle (go.libremail AT gmail DOT com)

package bits

import (
	"crypto/sha256"
	"crypto/rand"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"

	"libertymail-go/bits/ripemd160"
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
	if err != nil {
		b = nil
	}
	return b, err
}

/* 
AES decrypts data. This is done in-place, this function has a side effect.
*/
func AESDecrypt( data []byte, key []byte) ([]byte, error) {

	if len(data)%aes.BlockSize != 0 {
		return nil, errors.New("Data should be a multiple of 128 bits")
	}

	if len(data) < 32 {
		return nil, errors.New("iv + data should be at least 256 bits")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := data[:16]
	data = data[16:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(data, data)

	return data, nil
}

/*
AES encrypts data.
*/
func AESEncrypt(data []byte, key []byte) ([]byte, error) {

	if len(data)%16 != 0 {
		return nil, errors.New("Data should be a multiple of 128 bits")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, len(data) + 16)
	
	iv := ciphertext[:16]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[16:], data)

	return ciphertext, nil
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
