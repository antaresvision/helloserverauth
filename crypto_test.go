package main

import "testing"

func TestEncryptToBase64(t *testing.T) {
	txt := "Pippo"

	enc := EncryptToBase64(txt)
	t.Log(enc)
	dec := DecryptFromBase64(enc)

	if txt != dec {
		t.Fail()
	}

}

