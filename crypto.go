package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
)

//see https://godoc.org/golang.org/x/crypto/nacl/secretbox

var secretKey [32]byte

func init() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	secretKeyBytes, err := hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		panic(err)
	}

	copy(secretKey[:], secretKeyBytes)
}

func EncryptToBase64(clearText string) string {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// This encrypts clearText and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], []byte(clearText), &nonce, &secretKey)

	return base64.StdEncoding.EncodeToString(encrypted)
}

func DecryptFromBase64(encryptedText string) string {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		panic(err)
	}

	// When you decrypt, you must use the same nonce and key you used to
	// encrypt the message. One way to achieve this is to store the nonce
	// alongside the encrypted message. Above, we stored the nonce in the first
	// 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], encryptedBytes[:24])
	decrypted, ok := secretbox.Open(nil, encryptedBytes[24:], &decryptNonce, &secretKey)
	if !ok {
		panic("decryption error")
	}

	return string(decrypted)
}
