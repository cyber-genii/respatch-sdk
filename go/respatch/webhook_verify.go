package respatch

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DefaultWebhookToleranceSeconds is the default replay-protection window (5 minutes).
const DefaultWebhookToleranceSeconds int64 = 300

func decodeHexSHA256(s string) ([]byte, bool) {
	if len(s) != 64 {
		return nil, false
	}
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != sha256.Size {
		return nil, false
	}
	return b, true
}

// VerifyWebhookSignature checks X-Respatch-Signature (or legacy X-Transmit-Signature)
// using the default 300-second timestamp tolerance.
func VerifyWebhookSignature(payload []byte, signatureHeader, secret string) bool {
	return VerifyWebhookSignatureWithTolerance(payload, signatureHeader, secret, DefaultWebhookToleranceSeconds)
}

// VerifyWebhookSignatureWithTolerance checks X-Respatch-Signature (or legacy X-Transmit-Signature).
// It supports `t=<timestamp>,v1=<hmac_hex>` (including multiple `v1` entries for key rotation)
// as well as legacy `sha256=<hex>` and raw hex signatures.
// Pass toleranceSeconds <= 0 to bypass the timestamp age check (e.g. in unit tests).
func VerifyWebhookSignatureWithTolerance(payload []byte, signatureHeader, secret string, toleranceSeconds int64) bool {
	if signatureHeader == "" || secret == "" {
		return false
	}

	var timestamp string
	var v1Sigs []string
	for _, part := range strings.Split(signatureHeader, ",") {
		idx := strings.IndexByte(part, '=')
		if idx < 0 {
			continue
		}
		k := strings.TrimSpace(part[:idx])
		v := strings.TrimSpace(part[idx+1:])
		if k == "t" && timestamp == "" {
			timestamp = v
		} else if k == "v1" && v != "" {
			v1Sigs = append(v1Sigs, v)
		}
	}

	if timestamp != "" || len(v1Sigs) > 0 {
		if timestamp == "" || len(v1Sigs) == 0 {
			return false
		}
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || ts < 0 {
			return false
		}
		if toleranceSeconds > 0 {
			diff := time.Now().Unix() - ts
			if diff < 0 {
				diff = -diff
			}
			if diff > toleranceSeconds {
				return false
			}
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(timestamp))
		mac.Write([]byte("."))
		mac.Write(payload)
		expected := mac.Sum(nil)

		for _, v1 := range v1Sigs {
			provided, ok := decodeHexSHA256(v1)
			if ok && subtle.ConstantTimeCompare(expected, provided) == 1 {
				return true
			}
		}
		return false
	}

	providedStr := strings.TrimSpace(signatureHeader)
	if strings.HasPrefix(strings.ToLower(providedStr), "sha256=") {
		providedStr = strings.TrimSpace(providedStr[7:])
	}
	provided, ok := decodeHexSHA256(providedStr)
	if !ok {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)

	return subtle.ConstantTimeCompare(expected, provided) == 1
}

// SignWebhookPayload returns the hex digest for tests and replay tools.
func SignWebhookPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// SignWebhookHeader returns a `t=<timestamp>,v1=<hex>` signature header value.
func SignWebhookHeader(payload []byte, secret string, timestamp int64) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strconv.FormatInt(timestamp, 10)))
	mac.Write([]byte("."))
	mac.Write(payload)
	return fmt.Sprintf("t=%d,v1=%s", timestamp, hex.EncodeToString(mac.Sum(nil)))
}
