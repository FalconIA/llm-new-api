package signature

import (
	"crypto/aes"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCMCShanghaiSign(t *testing.T) {
	const key = "1234567890abcdef"
	now := time.UnixMilli(1700000000000)
	signed, err := signAt(key, "app_1", 2, now)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(signed, "app_1."))

	encoded := strings.TrimPrefix(signed, "app_1.")
	encrypted, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	block, err := aes.NewCipher([]byte(key))
	require.NoError(t, err)
	plaintext := make([]byte, len(encrypted))
	for i := 0; i < len(encrypted); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], encrypted[i:i+aes.BlockSize])
	}
	padding := int(plaintext[len(plaintext)-1])
	require.GreaterOrEqual(t, padding, 1)
	require.LessOrEqual(t, padding, aes.BlockSize)
	assert.Equal(t, `{"time":1700000002000}`, string(plaintext[:len(plaintext)-padding]))
}

func TestCMCShanghaiRejectsInvalidConfiguration(t *testing.T) {
	assert.Error(t, Validate("short", "app", 0))
	assert.Error(t, Validate("1234567890abcdef", "", 0))
	assert.Error(t, Validate("1234567890abcdef", "app\tname", 0))
	assert.Error(t, Validate("1234567890abcdef", "app", 3601))
}
