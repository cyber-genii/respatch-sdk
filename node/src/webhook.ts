import { createHmac, timingSafeEqual } from 'crypto';

export const RESPATCH_SIGNATURE_HEADER = 'X-Respatch-Signature';

/** @deprecated Use RESPATCH_SIGNATURE_HEADER */
export const TRANSMIT_SIGNATURE_HEADER = RESPATCH_SIGNATURE_HEADER;

/**
 * Verify an inbound Respatch webhook using `X-Respatch-Signature` (or legacy `X-Transmit-Signature`).
 * Signature is HMAC-SHA256 of the raw request body, hex-encoded.
 */
export function verifyWebhookSignature(
  payload: string | Buffer,
  signatureHeader: string,
  secret: string,
): boolean {
  if (!signatureHeader || !secret) return false;

  const body = typeof payload === 'string' ? payload : payload.toString('utf8');

  // Support standard t={timestamp},v1={hmac} signature format
  if (signatureHeader.includes('t=') && signatureHeader.includes('v1=')) {
    const parts = signatureHeader.split(',');
    let timestamp = '';
    let v1 = '';
    for (const p of parts) {
      const idx = p.indexOf('=');
      if (idx !== -1) {
        const k = p.slice(0, idx).trim();
        const v = p.slice(idx + 1).trim();
        if (k === 't') timestamp = v;
        if (k === 'v1') v1 = v;
      }
    }
    if (timestamp && v1) {
      const expected = createHmac('sha256', secret)
        .update(`${timestamp}.${body}`, 'utf8')
        .digest('hex');
      try {
        const a = Buffer.from(expected, 'hex');
        const b = Buffer.from(v1, 'hex');
        if (a.length === b.length && timingSafeEqual(a, b)) {
          return true;
        }
      } catch {
        // Fall back to plain signature check
      }
    }
  }

  const expected = createHmac('sha256', secret).update(body, 'utf8').digest('hex');
  const provided = signatureHeader.replace(/^sha256=/i, '').trim();

  try {
    const a = Buffer.from(expected, 'hex');
    const b = Buffer.from(provided, 'hex');
    return a.length === b.length && timingSafeEqual(a, b);
  } catch {
    return false;
  }
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
