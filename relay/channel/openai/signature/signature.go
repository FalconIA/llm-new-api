package signature

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// CMCShanghai is the legacy provider protocol: AES ECB over a padded
// {"time":<Unix milliseconds>} payload, prefixed by the application ID.
const CMCShanghai = "cmc_sh"

func Validate(key, appID string, offsetSeconds int64) error {
	switch len(key) {
	case 16, 24, 32:
	default:
		return errors.New("cmc_sh key must be 16, 24, or 32 bytes")
	}
	if len(appID) == 0 || len(appID) > 128 || strings.IndexFunc(appID, unicode.IsSpace) >= 0 {
		return errors.New("signature_app_id must be 1-128 characters without whitespace")
	}
	if offsetSeconds < -3600 || offsetSeconds > 3600 {
		return errors.New("signature_time_offset must be within one hour")
	}
	return nil
}

func Sign(key, appID string, offsetSeconds int64) (string, error) {
	return signAt(key, appID, offsetSeconds, time.Now())
}

func signAt(key, appID string, offsetSeconds int64, now time.Time) (string, error) {
	if err := Validate(key, appID, offsetSeconds); err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("create cmc_sh cipher: %w", err)
	}
	payload := []byte(`{"time":` + strconv.FormatInt(now.UnixMilli()+offsetSeconds*1000, 10) + `}`)
	padding := aes.BlockSize - len(payload)%aes.BlockSize
	for range padding {
		payload = append(payload, byte(padding))
	}
	encrypted := make([]byte, len(payload))
	for i := 0; i < len(payload); i += aes.BlockSize {
		block.Encrypt(encrypted[i:i+aes.BlockSize], payload[i:i+aes.BlockSize])
	}
	return appID + "." + base64.StdEncoding.EncodeToString(encrypted), nil
}
