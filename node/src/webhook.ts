import { createHmac, timingSafeEqual } from 'crypto';

export const RESPATCH_SIGNATURE_HEADER = 'X-Respatch-Signature';
export const DEFAULT_WEBHOOK_TOLERANCE_SECONDS = 300;

/** @deprecated Use RESPATCH_SIGNATURE_HEADER */
export const TRANSMIT_SIGNATURE_HEADER = RESPATCH_SIGNATURE_HEADER;

const HEX_SHA256_RE = /^[0-9a-fA-F]{64}$/;

function parseHexSha256(hex: string): Buffer | null {
  if (!HEX_SHA256_RE.test(hex)) return null;
  return Buffer.from(hex, 'hex');
}

/**
 * Verify an inbound Respatch webhook using `X-Respatch-Signature` (or legacy `X-Transmit-Signature`).
 * Supports `t=<timestamp>,v1=<hmac_hex>` (including multiple `v1` entries for key rotation)
 * with replay window tolerance, as well as legacy `sha256=<hex>` / raw hex signatures.
 *
 * Pass `toleranceSeconds <= 0` to bypass the timestamp age check (e.g. in unit tests).
 */
export function verifyWebhookSignature(
  payload: string | Buffer,
  signatureHeader: string,
  secret: string,
  toleranceSeconds: number = DEFAULT_WEBHOOK_TOLERANCE_SECONDS,
): boolean {
  if (!signatureHeader || !secret) return false;

  const rawBody = typeof payload === 'string' ? Buffer.from(payload, 'utf8') : payload;

  // Parse structured t=<timestamp>,v1=<hmac>[,v1=<hmac2>...] format
  let timestamp = '';
  const v1Signatures: string[] = [];
  for (const part of signatureHeader.split(',')) {
    const idx = part.indexOf('=');
    if (idx === -1) continue;
    const k = part.slice(0, idx).trim();
    const v = part.slice(idx + 1).trim();
    if (k === 't' && !timestamp) {
      timestamp = v;
    } else if (k === 'v1' && v) {
      v1Signatures.push(v);
    }
  }

  if (timestamp || v1Signatures.length > 0) {
    if (!timestamp || v1Signatures.length === 0 || !/^\d+$/.test(timestamp)) {
      return false;
    }
    const ts = Number(timestamp);
    if (!Number.isSafeInteger(ts)) {
      return false;
    }
    if (toleranceSeconds > 0) {
      const nowSeconds = Math.floor(Date.now() / 1000);
      if (Math.abs(nowSeconds - ts) > toleranceSeconds) {
        return false;
      }
    }

    const expected = createHmac('sha256', secret)
      .update(Buffer.from(`${timestamp}.`, 'utf8'))
      .update(rawBody)
      .digest();

    for (const v1 of v1Signatures) {
      const provided = parseHexSha256(v1);
      if (provided && provided.length === expected.length && timingSafeEqual(expected, provided)) {
        return true;
      }
    }
    return false;
  }

  const expected = createHmac('sha256', secret).update(rawBody).digest();
  const provided = parseHexSha256(signatureHeader.replace(/^sha256=/i, '').trim());
  if (!provided) return false;
  return provided.length === expected.length && timingSafeEqual(expected, provided);
}

/** Resolve signature header from incoming request headers (Respatch preferred). */
export function getWebhookSignatureHeader(
  headers: Record<string, string | string[] | undefined>,
): string {
  const respatch = headers[RESPATCH_SIGNATURE_HEADER.toLowerCase()];
  if (typeof respatch === 'string') return respatch;
  const legacy = headers['x-transmit-signature'];
  if (typeof legacy === 'string') return legacy;
  return '';
}

export function signWebhookPayload(payload: string | Buffer, secret: string): string {
  const body = typeof payload === 'string' ? payload : payload.toString('utf8');
  return createHmac('sha256', secret).update(body, 'utf8').digest('hex');
}

/** Sign a payload in the `t=<timestamp>,v1=<hmac>` header format. */
export function signWebhookHeader(
  payload: string | Buffer,
  secret: string,
  timestamp: number = Math.floor(Date.now() / 1000),
): string {
  const rawBody = typeof payload === 'string' ? Buffer.from(payload, 'utf8') : payload;
  const sig = createHmac('sha256', secret)
    .update(Buffer.from(`${timestamp}.`, 'utf8'))
    .update(rawBody)
    .digest('hex');
  return `t=${timestamp},v1=${sig}`;
}
