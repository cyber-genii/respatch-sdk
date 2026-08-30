package respatch

import "testing"

func TestWebhookSignatureRoundTrip(t *testing.T) {
	secret := "whsec_test"
	body := []byte(`{"event":"delivery.status_changed"}`)
	sig := SignWebhookPayload(body, secret)
	if !VerifyWebhookSignature(body, sig, secret) {
		t.Fatal("expected valid signature")
	}
	if !VerifyWebhookSignature(body, "sha256="+sig, secret) {
		t.Fatal("expected valid sha256= prefixed signature")
	}
	if VerifyWebhookSignature(append(body, 'x'), sig, secret) {
		t.Fatal("expected tampered payload to fail")
	}
}
