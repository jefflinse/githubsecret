package githubsecret

import (
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockValidSecretValue          = "a valid secret value"
	mockValidPublicKey            = "hel9i9lSE4Cu103BBICvKhmLi8LLnVO7BDdqANPOlEw="
	mockInvalidPublicKeyTooShort  = "z8Viu/+IGVACPAltd3UpMCBWV+yxUZkDXkQcQwK/"
	mockInvalidPublicKeyTooLong   = "mSu6F6vPOCU7inMJ1CNsXPCc1f/oN6hgfSWjuxvoxQTeeA=="
	mockInvalidPublicKeyNotBase64 = "this is certainly not base64!"

	mockGeneratedPublicKey  = "Gj5IXxUxMQ5tEVN7a83nwiBC/e5ckiWTkPRH3G3FNBw="
	mockGeneratedPrivateKey = "uXts8kET12LelZEmqhnscs/0Qj4CpC3V3yWhFfo9czk="
)

func TestEncrypt(t *testing.T) {
	for _, test := range []struct {
		name           string
		pk             string
		secret         string
		generateKeyFn  func(r io.Reader) (*[32]byte, *[32]byte, error)
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "can encode an empty secret with a valid recipient pk",
			pk:             mockValidPublicKey,
			secret:         "",
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "Gj5IXxUxMQ5tEVN7a83nwiBC/e5ckiWTkPRH3G3FNBxAvNLVFmltuTJxCdc5joBu",
			expectedError:  nil,
		},
		{
			name:           "can encode a non-empty secret with a valid recipient pk",
			pk:             mockValidPublicKey,
			secret:         mockValidSecretValue,
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "Gj5IXxUxMQ5tEVN7a83nwiBC/e5ckiWTkPRH3G3FNBzNnh9ZMlNTPy2fTY6WSlvJ3CBqQVYj2jQwqwMSTcShAiaQCIk=",
			expectedError:  nil,
		},
		{
			name:           "fails when the recipient pk is empty",
			pk:             "",
			secret:         mockValidSecretValue,
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "",
			expectedError:  errors.New("recipient public key has invalid length (0 bytes)"),
		},
		{
			name:           "fails when the recipient pk is too short",
			pk:             mockInvalidPublicKeyTooShort,
			secret:         mockValidSecretValue,
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "",
			expectedError:  errors.New("recipient public key has invalid length (30 bytes)"),
		},
		{
			name:           "fails when the recipient pk is too long",
			pk:             mockInvalidPublicKeyTooLong,
			secret:         mockValidSecretValue,
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "",
			expectedError:  errors.New("recipient public key has invalid length (34 bytes)"),
		},
		{
			name:           "fails when the recipient pk is not valid base64",
			pk:             mockInvalidPublicKeyNotBase64,
			secret:         mockValidSecretValue,
			generateKeyFn:  mockGenerateKey,
			expectedOutput: "",
			expectedError:  errors.New("illegal base64 data at input byte 4"),
		},
		{
			name:   "fails when the ephemeral key pair cannot be generated",
			pk:     mockValidPublicKey,
			secret: mockValidSecretValue,
			generateKeyFn: func(_ io.Reader) (*[32]byte, *[32]byte, error) {
				return nil, nil, assert.AnError
			},
			expectedOutput: "",
			expectedError:  assert.AnError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			origGenerateKey := generateKey
			defer func() { generateKey = origGenerateKey }()
			generateKey = test.generateKeyFn

			output, err := Encrypt(test.pk, test.secret)
			if test.expectedError != nil {
				assert.EqualError(t, err, test.expectedError.Error())
				assert.Empty(t, output)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, output)
				assert.True(t, strings.HasPrefix(output, mockGeneratedPublicKey[:len(mockGeneratedPublicKey)-2]))
			}
		})
	}
}

// Returns a predefined keypair for mocking out randomness in nonce generation.
func mockGenerateKey(r io.Reader) (*[32]byte, *[32]byte, error) {
	pub := new([keySize]byte)
	base64.StdEncoding.Decode(pub[:], []byte(mockGeneratedPublicKey))
	priv := new([keySize]byte)
	base64.StdEncoding.Decode(priv[:], []byte(mockGeneratedPrivateKey))
	return pub, priv, nil
}
