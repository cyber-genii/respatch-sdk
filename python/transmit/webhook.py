"""Webhook signature verification for inbound Respatch events."""

from __future__ import annotations

import hashlib
import hmac
import re
import time
from typing import Optional

RESPATCH_SIGNATURE_HEADER = "X-Respatch-Signature"
DEFAULT_WEBHOOK_TOLERANCE_SECONDS = 300
_HEX_SHA256_RE = re.compile(r"^[0-9a-fA-F]{64}$")


def verify_webhook_signature(
    payload: str | bytes,
    signature_header: str,
    secret: str,
    tolerance_seconds: int = DEFAULT_WEBHOOK_TOLERANCE_SECONDS,
) -> bool:
    """Verify ``X-Respatch-Signature`` (or legacy ``X-Transmit-Signature``).

    Supports ``t=<timestamp>,v1=<hmac_hex>`` (including multiple ``v1`` entries
    for key rotation) with replay window tolerance, as well as legacy
    ``sha256=<hex>`` and raw hex signatures.

    Pass ``tolerance_seconds <= 0`` to bypass the timestamp age check.
    """
    if not signature_header or not secret:
        return False

    body = payload.encode("utf-8") if isinstance(payload, str) else payload
    secret_bytes = secret.encode("utf-8")

    # Parse structured t=<timestamp>,v1=<hmac>[,v1=<hmac2>...] format
    timestamp: Optional[str] = None
    v1_signatures: list[str] = []
    for part in signature_header.split(","):
        if "=" not in part:
            continue
        k, v = part.split("=", 1)
        k = k.strip()
        v = v.strip()
        if k == "t" and timestamp is None:
            timestamp = v
        elif k == "v1" and v:
            v1_signatures.append(v)

    if timestamp is not None or v1_signatures:
        if not timestamp or not v1_signatures or not timestamp.isdigit():
            return False
        ts = int(timestamp)
        if tolerance_seconds > 0:
            if abs(int(time.time()) - ts) > tolerance_seconds:
                return False

        signed_payload = f"{timestamp}.".encode("utf-8") + body
        expected = hmac.new(secret_bytes, signed_payload, hashlib.sha256).hexdigest()
        for v1 in v1_signatures:
            if _HEX_SHA256_RE.match(v1) and hmac.compare_digest(expected, v1.lower()):
                return True
        return False

    provided = signature_header.strip()
    if provided.lower().startswith("sha256="):
        provided = provided[7:].strip()
    if not _HEX_SHA256_RE.match(provided):
        return False
    expected = hmac.new(secret_bytes, body, hashlib.sha256).hexdigest()
    return hmac.compare_digest(expected, provided.lower())


def get_webhook_signature_header(headers: dict) -> str:
    for key in (RESPATCH_SIGNATURE_HEADER, "X-Transmit-Signature"):
        val = headers.get(key) or headers.get(key.lower())
        if isinstance(val, str) and val:
            return val
    return ""


def sign_webhook_payload(payload: str | bytes, secret: str) -> str:
    body = payload.encode("utf-8") if isinstance(payload, str) else payload
    return hmac.new(secret.encode("utf-8"), body, hashlib.sha256).hexdigest()


def sign_webhook_header(
    payload: str | bytes,
    secret: str,
    timestamp: Optional[int] = None,
) -> str:
    """Sign a payload in the ``t=<timestamp>,v1=<hmac>`` header format."""
    ts = int(time.time()) if timestamp is None else int(timestamp)
    body = payload.encode("utf-8") if isinstance(payload, str) else payload
    signed_payload = f"{ts}.".encode("utf-8") + body
    sig = hmac.new(secret.encode("utf-8"), signed_payload, hashlib.sha256).hexdigest()
    return f"t={ts},v1={sig}"
