package respatch

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestWebhookSignatureRoundTrip(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"event":"delivery.status_changed"}`)
	sig := SignWebhookPayload(body, secret)
	if !VerifyWebhookSignature(body, sig, secret) {
		t.Fatal("expected valid legacy signature")
	}
	if !VerifyWebhookSignature(body, "sha256="+sig, secret) {
		t.Fatal("expected valid sha256= prefixed signature")
	}
	if VerifyWebhookSignature(append([]byte{}, append(body, 'x')...), sig, secret) {
		t.Fatal("expected tampered payload to fail")
	}
}

func TestWebhookSignatureTimestampedAndRotation(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"event":"delivery.status_changed","data":{"id":"1"}}`)
	now := time.Now().Unix()

	header := SignWebhookHeader(body, secret, now)
	if !VerifyWebhookSignature(body, header, secret) {
		t.Fatal("expected valid t=<ts>,v1=<sig> signature")
	}

	parts := strings.SplitN(header, "v1=", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected header format: %s", header)
	}
	rotated := fmt.Sprintf("t=%d,v1=%s,v1=%s", now, strings.Repeat("0", 64), parts[1])
	if !VerifyWebhookSignature(body, rotated, secret) {
		t.Fatal("expected rotated multiple v1 header to verify")
	}

	stale := SignWebhookHeader(body, secret, now-600)
	if VerifyWebhookSignature(body, stale, secret) {
		t.Fatal("expected stale timestamp to fail default tolerance")
	}
	if !VerifyWebhookSignatureWithTolerance(body, stale, secret, 0) {
		t.Fatal("expected tolerance <= 0 to bypass timestamp check")
	}

	future := SignWebhookHeader(body, secret, now+600)
	if VerifyWebhookSignature(body, future, secret) {
		t.Fatal("expected future timestamp beyond tolerance to fail")
	}

	if VerifyWebhookSignature(body, header, "wrong") {
		t.Fatal("expected wrong secret to fail")
	}
	if VerifyWebhookSignature(append([]byte{}, append(body, 'x')...), header, secret) {
		t.Fatal("expected tampered payload to fail")
	}
	if VerifyWebhookSignature(body, "t=bad,v1="+strings.Repeat("a", 64), secret) {
		t.Fatal("expected invalid timestamp to fail")
	}
	if VerifyWebhookSignature(body, "t=12345", secret) {
		t.Fatal("expected missing v1 to fail")
	}
}
