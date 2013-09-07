package address

import (
	"testing"
)

func TestAddress(t *testing.T) {

	addr, err := NewAddress(1, 0)
	if err != nil {
		t.Error(err.Error())
	}

	if !ValidateIdentifier(addr.Identifier) {
		t.Error("Invalid address identifier", addr.Identifier)
	}

	if !ValidateChecksum(addr.Identifier) {
		t.Error("Invalid checksum")
	}
}
