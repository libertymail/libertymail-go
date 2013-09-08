package hashing

import (
	"encoding/base64"
	"testing"
)

func TestSHA(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// SHA512(This is a string.)
	//
	// hex: 0145c77435b886e43fcfa5b8a6e2c5a9f1c216f694a65e75354f9679174551b7a0151b72f5497d58845bc5033f39f3249ee087cdb602680edc3fdeda8a18ff9b
	// HEX: 0145C77435B886E43FCFA5B8A6E2C5A9F1C216F694A65E75354F9679174551B7A0151B72F5497D58845BC5033F39F3249EE087CDB602680EDC3FDEDA8A18FF9B
	// base64: AUXHdDW4huQ/z6W4puLFqfHCFvaUpl51NU+WeRdFUbegFRty9Ul9WIRbxQM/OfMknuCHzbYCaA7cP97aihj/mw==

	msg := "This is a string."

	digest := SHA512([]byte(msg))
	encoded := base64.StdEncoding.EncodeToString(digest)

	if encoded != "AUXHdDW4huQ/z6W4puLFqfHCFvaUpl51NU+WeRdFUbegFRty9Ul9WIRbxQM/OfMknuCHzbYCaA7cP97aihj/mw==" {
		t.Error("Sha512 did not match")
	}
}
