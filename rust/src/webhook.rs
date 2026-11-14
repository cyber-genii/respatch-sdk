use hmac::{Hmac, Mac};
use sha2::Sha256;
use std::time::{SystemTime, UNIX_EPOCH};

type HmacSha256 = Hmac<Sha256>;

pub const DEFAULT_WEBHOOK_TOLERANCE_SECONDS: i64 = 300;

fn decode_hex_sha256(s: &str) -> Option<[u8; 32]> {
    if s.len() != 64 {
        return None;
    }
    let bytes = hex::decode(s).ok()?;
    let arr: [u8; 32] = bytes.try_into().ok()?;
    Some(arr)
}

/// Verify `X-Respatch-Signature` (or legacy `X-Transmit-Signature`) using the default
/// 300-second timestamp tolerance.
pub fn verify_webhook_signature(payload: &[u8], signature_header: &str, secret: &str) -> bool {
    verify_webhook_signature_with_tolerance(
        payload,
        signature_header,
        secret,
        DEFAULT_WEBHOOK_TOLERANCE_SECONDS,
    )
}

/// Verify `X-Respatch-Signature` (or legacy `X-Transmit-Signature`) with a configurable
/// timestamp age tolerance in seconds.
///
/// Supports `t=<timestamp>,v1=<hmac_hex>` (including multiple `v1` entries for key rotation),
/// as well as legacy `sha256=<hex>` and raw hex signatures.
/// Pass `tolerance_seconds <= 0` to bypass the timestamp age check (e.g. in unit tests).
pub fn verify_webhook_signature_with_tolerance(
    payload: &[u8],
    signature_header: &str,
    secret: &str,
    tolerance_seconds: i64,
) -> bool {
    if signature_header.is_empty() || secret.is_empty() {
        return false;
    }

    let mut timestamp: Option<&str> = None;
    let mut v1_sigs: Vec<&str> = Vec::new();
    for part in signature_header.split(',') {
        if let Some((k, v)) = part.split_once('=') {
            let k = k.trim();
            let v = v.trim();
            if k == "t" && timestamp.is_none() {
                timestamp = Some(v);
            } else if k == "v1" && !v.is_empty() {
                v1_sigs.push(v);
            }
        }
    }

    if timestamp.is_some() || !v1_sigs.is_empty() {
        let Some(ts_str) = timestamp else {
            return false;
        };
        if v1_sigs.is_empty() || !ts_str.bytes().all(|b| b.is_ascii_digit()) {
            return false;
        }
        let Ok(ts) = ts_str.parse::<i64>() else {
            return false;
        };
        if tolerance_seconds > 0 {
            let Ok(now) = SystemTime::now().duration_since(UNIX_EPOCH) else {
                return false;
            };
            let now_secs = now.as_secs() as i64;
            if (now_secs - ts).abs() > tolerance_seconds {
                return false;
            }
        }

        let mut mac = match HmacSha256::new_from_slice(secret.as_bytes()) {
            Ok(m) => m,
            Err(_) => return false,
        };
        mac.update(ts_str.as_bytes());
        mac.update(b".");
        mac.update(payload);
        let expected = mac.finalize().into_bytes();

        return v1_sigs.into_iter().any(|v1| {
            decode_hex_sha256(v1)
                .map(|provided| constant_time_eq(&expected, &provided))
                .unwrap_or(false)
        });
    }

    let trimmed = signature_header.trim();
    let provided_hex = if trimmed.len() >= 7 && trimmed[..7].eq_ignore_ascii_case("sha256=") {
        trimmed[7..].trim()
    } else {
        trimmed
    };
    let Some(provided) = decode_hex_sha256(provided_hex) else {
        return false;
    };

    let mut mac = match HmacSha256::new_from_slice(secret.as_bytes()) {
        Ok(m) => m,
        Err(_) => return false,
    };
    mac.update(payload);
    let expected = mac.finalize().into_bytes();
    constant_time_eq(&expected, &provided)
}

/// Sign a payload for tests or replay tools (raw hex digest).
pub fn sign_webhook_payload(payload: &[u8], secret: &str) -> String {
    let mut mac = HmacSha256::new_from_slice(secret.as_bytes()).expect("hmac key");
    mac.update(payload);
    hex::encode(mac.finalize().into_bytes())
}

/// Sign a payload in the `t=<timestamp>,v1=<hmac_hex>` header format.
pub fn sign_webhook_header(payload: &[u8], secret: &str, timestamp: i64) -> String {
    let ts_str = timestamp.to_string();
    let mut mac = HmacSha256::new_from_slice(secret.as_bytes()).expect("hmac key");
    mac.update(ts_str.as_bytes());
    mac.update(b".");
    mac.update(payload);
    format!("t={},v1={}", ts_str, hex::encode(mac.finalize().into_bytes()))
}

fn constant_time_eq(a: &[u8], b: &[u8]) -> bool {
    if a.len() != b.len() {
        return false;
    }
    a.iter()
        .zip(b.iter())
        .fold(0u8, |acc, (x, y)| acc | (x ^ y))
        == 0
}

#[cfg(test)]
mod tests {
    use super::*;

    fn now_secs() -> i64 {
        SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs() as i64
    }

    #[test]
    fn round_trip_legacy_signature() {
        let secret = "whsec_test";
        let body = br#"{"event":"delivery.status_changed"}"#;
        let sig = sign_webhook_payload(body, secret);
        assert!(verify_webhook_signature(body, &sig, secret));
        assert!(verify_webhook_signature(body, &format!("sha256={sig}"), secret));
        assert!(!verify_webhook_signature(br#"{"event":"tampered"}"#, &sig, secret));
    }

    #[test]
    fn round_trip_timestamped_and_rotation() {
        let secret = "whsec_test";
        let body = br#"{"event":"delivery.status_changed","data":{"id":"1"}}"#;
        let now = now_secs();
        let header = sign_webhook_header(body, secret, now);
        assert!(verify_webhook_signature(body, &header, secret));

        let v1 = header.split("v1=").nth(1).unwrap();
        let rotated = format!("t={now},v1={},v1={v1}", "0".repeat(64));
        assert!(verify_webhook_signature(body, &rotated, secret));

        let stale = sign_webhook_header(body, secret, now - 600);
        assert!(!verify_webhook_signature(body, &stale, secret));
        assert!(verify_webhook_signature_with_tolerance(body, &stale, secret, 0));
        assert!(verify_webhook_signature_with_tolerance(body, &stale, secret, -1));

        let future = sign_webhook_header(body, secret, now + 600);
        assert!(!verify_webhook_signature(body, &future, secret));

        assert!(!verify_webhook_signature(body, &header, "wrong"));
        assert!(!verify_webhook_signature(br#"{"tampered":true}"#, &header, secret));
        assert!(!verify_webhook_signature(body, &format!("t=bad,v1={}", "a".repeat(64)), secret));
        assert!(!verify_webhook_signature(body, "t=12345", secret));
    }
}
