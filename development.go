package main

import (
	"fmt"
	"bytes"

    "libertymail-go/bits"
)

func main() {
    fmt.Println("Hello, world!")
    key := []byte("example key 1234")
    plaintext := []byte("exampleplaintext")

	fmt.Println(plaintext)
    enc, _ := bits.AESEncrypt(plaintext, key)
    fmt.Println(enc)
    dec, _ := bits.AESDecrypt(enc, key)

    fmt.Println(dec)
    fmt.Println(bytes.Equal(plaintext,dec))
}



