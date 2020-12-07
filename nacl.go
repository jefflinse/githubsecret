package githubsecret

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

const (
	keySize   = 32
	nonceSize = 24
)

// Encrypt encrypts a secret using the provided recipient public key.
func Encrypt(recipientPublicKey string, content string) (string, error) {
	// decode the provided public key from base64
	recipientKey := new([keySize]byte)
	b, err := base64.StdEncoding.DecodeString(recipientPublicKey)
	if err != nil {
		return "", err
	}

	copy(recipientKey[:], b)

	// create an ephemeral key pair
	pubKey, privKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}

	// create the nonce by hashing together the two public keys
	nonce := new([nonceSize]byte)
	nonceHash, err := blake2b.New(nonceSize, nil)
	if err != nil {
		return "", err
	}

	if _, err := nonceHash.Write(pubKey[:]); err != nil {
		return "", err
	}

	if _, err := nonceHash.Write(recipientKey[:]); err != nil {
		return "", err
	}

	copy(nonce[:], nonceHash.Sum([]byte{}))

	// begin the output with the ephemeral public key and append the encrypted content
	out := box.Seal(pubKey[:], []byte(content), nonce, recipientKey, privKey)

	// base64-encode the final output
	return base64.StdEncoding.EncodeToString(out), nil
}
