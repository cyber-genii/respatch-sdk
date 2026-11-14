import sys
import time
import types
import unittest

try:
    import requests  # noqa: F401
except ModuleNotFoundError:
    _stub = types.ModuleType("requests")
    _stub.Session = object
    _stub.RequestException = Exception
    sys.modules["requests"] = _stub

import respatch
from respatch import (
    sign_webhook_header,
    sign_webhook_payload,
    verify_webhook_signature,
)


class WebhookSignatureTests(unittest.TestCase):
    secret = "whsec_test_secret"
    payload = b'{"event":"delivery.status_changed","data":{"id":"1"}}'

    def test_respatch_package_exports(self):
        self.assertIs(respatch.verify_webhook_signature, verify_webhook_signature)
        self.assertIs(respatch.sign_webhook_payload, sign_webhook_payload)
        self.assertIs(respatch.sign_webhook_header, sign_webhook_header)

    def test_legacy_signature(self):
        sig = sign_webhook_payload(self.payload, self.secret)
        self.assertTrue(verify_webhook_signature(self.payload, sig, self.secret))
        self.assertTrue(verify_webhook_signature(self.payload, f"sha256={sig}", self.secret))

    def test_timestamped_v1_signature(self):
        now = int(time.time())
        header = sign_webhook_header(self.payload, self.secret, timestamp=now)
        self.assertTrue(verify_webhook_signature(self.payload, header, self.secret))
        self.assertTrue(
            verify_webhook_signature(self.payload.decode("utf-8"), header, self.secret)
        )

    def test_multiple_v1_signatures_rotation(self):
        now = int(time.time())
        header = sign_webhook_header(self.payload, self.secret, timestamp=now)
        v1 = header.split("v1=")[1]
        rotated = f"t={now},v1={'0' * 64},v1={v1}"
        self.assertTrue(verify_webhook_signature(self.payload, rotated, self.secret))

    def test_timestamp_tolerance_and_bypass(self):
        stale = int(time.time()) - 600
        header = sign_webhook_header(self.payload, self.secret, timestamp=stale)
        self.assertFalse(verify_webhook_signature(self.payload, header, self.secret))
        self.assertTrue(
            verify_webhook_signature(self.payload, header, self.secret, tolerance_seconds=0)
        )
        self.assertTrue(
            verify_webhook_signature(self.payload, header, self.secret, tolerance_seconds=-1)
        )

        future = int(time.time()) + 600
        future_header = sign_webhook_header(self.payload, self.secret, timestamp=future)
        self.assertFalse(verify_webhook_signature(self.payload, future_header, self.secret))

    def test_rejects_tampered_and_malformed(self):
        header = sign_webhook_header(self.payload, self.secret)
        self.assertFalse(verify_webhook_signature(self.payload + b"x", header, self.secret))
        self.assertFalse(verify_webhook_signature(self.payload, header, "wrong"))
        self.assertFalse(
            verify_webhook_signature(self.payload, f"t=bad,v1={'a' * 64}", self.secret)
        )
        self.assertFalse(verify_webhook_signature(self.payload, "t=12345", self.secret))
        self.assertFalse(verify_webhook_signature(self.payload, "", self.secret))
        self.assertFalse(verify_webhook_signature(self.payload, header, ""))


if __name__ == "__main__":
    unittest.main()
