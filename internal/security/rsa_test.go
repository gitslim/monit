package security_test

import (
	"crypto/rsa"
	"testing"

	"github.com/gitslim/monit/internal/security"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptMessage(t *testing.T) {
	publicKey, err := security.ReadRSAPublicKeyFromFile("../../testdata/keys/public.pem")
	assert.NoError(t, err, "Failed to read public key")
	privateKey, err := security.ReadRSAPrivateKeyFromFile("../../testdata/keys/private.pem")
	assert.NoError(t, err, "Failed to read private key")

	tests := []struct {
		name       string
		message    []byte
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
		wantErr    bool
	}{
		{
			name:       "Simple message",
			message:    []byte("Simple message"),
			publicKey:  publicKey,
			privateKey: privateKey,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := security.EncryptRSA(tt.publicKey, tt.message)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("EncryptDecryptMessage() failed: %v", err)
				}
				return
			}
			decrypted, err := security.DecryptRSA(tt.privateKey, encrypted)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("EncryptDecryptMessage() failed: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EncryptDecryptMessage() succeeded unexpectedly")
			}

			assert.Equal(t, tt.message, decrypted)
		})
	}
}
