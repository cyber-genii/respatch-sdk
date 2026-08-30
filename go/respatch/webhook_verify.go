package respatch

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
)

// VerifyWebhookSignature checks X-Respatch-Signature (or legacy X-Transmit-Signature).
func VerifyWebhookSignature(payload []byte, signatureHeader, secret string) bool {
	if signatureHeader == "" || secret == "" {
		return false
	}

	provided := strings.TrimSpace(signatureHeader)
	if strings.HasPrefix(strings.ToLower(provided), "sha256=") {
		provided = strings.TrimSpace(provided[7:])
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	return subtle.ConstantTimeCompare([]byte(expected), []byte(provided)) == 1
}

// SignWebhookPayload returns the hex digest for tests and replay tools.
func SignWebhookPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
