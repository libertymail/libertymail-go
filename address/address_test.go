package address

import (
	"testing"
)

func TestAddress(t *testing.T) {

	addr, err := NewAddress(1, 0, false)
	if err != nil {
		t.Error(err.Error())
	}

	if !ValidateIdentifier(addr.Identifier) {
		t.Error("Invalid address identifier", addr.Identifier)
	}

	valid, err := ValidateChecksum(addr.Identifier)
	if err != nil {
		t.Error(err.Error())
	}

	if valid != true {
		t.Error("Invalid checksum")
	}
}
