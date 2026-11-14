import { expect } from 'chai';
import {
  signWebhookHeader,
  signWebhookPayload,
  verifyWebhookSignature,
} from '../src/webhook';

describe('webhook signatures', () => {
  const secret = 'whsec_test_secret';
  const payload = JSON.stringify({ event: 'delivery.status_changed', data: { id: '1' } });

  it('verifies a valid legacy signature', () => {
    const sig = signWebhookPayload(payload, secret);
    expect(verifyWebhookSignature(payload, sig, secret)).to.equal(true);
    expect(verifyWebhookSignature(payload, `sha256=${sig}`, secret)).to.equal(true);
  });

  it('verifies a valid t=<ts>,v1=<sig> signature within tolerance', () => {
    const now = Math.floor(Date.now() / 1000);
    const header = signWebhookHeader(payload, secret, now);
    expect(verifyWebhookSignature(payload, header, secret)).to.equal(true);
  });

  it('verifies when multiple v1 signatures are present (secret rotation)', () => {
    const now = Math.floor(Date.now() / 1000);
    const valid = signWebhookHeader(payload, secret, now);
    const v1 = valid.split('v1=')[1];
    const rotatedHeader = `t=${now},v1=${'0'.repeat(64)},v1=${v1}`;
    expect(verifyWebhookSignature(payload, rotatedHeader, secret)).to.equal(true);
  });

  it('rejects expired timestamps and allows tolerance <= 0 bypass', () => {
    const stale = Math.floor(Date.now() / 1000) - 600;
    const header = signWebhookHeader(payload, secret, stale);
    expect(verifyWebhookSignature(payload, header, secret)).to.equal(false);
    expect(verifyWebhookSignature(payload, header, secret, 0)).to.equal(true);
    expect(verifyWebhookSignature(payload, header, secret, -1)).to.equal(true);
  });

  it('rejects future timestamps beyond tolerance', () => {
    const future = Math.floor(Date.now() / 1000) + 600;
    const header = signWebhookHeader(payload, secret, future);
    expect(verifyWebhookSignature(payload, header, secret)).to.equal(false);
  });

  it('rejects tampered payloads and malformed headers', () => {
    const sig = signWebhookPayload(payload, secret);
    expect(verifyWebhookSignature(`${payload}x`, sig, secret)).to.equal(false);

    const header = signWebhookHeader(payload, secret);
    expect(verifyWebhookSignature(`${payload}x`, header, secret)).to.equal(false);
    expect(verifyWebhookSignature(payload, 't=notanumber,v1=' + 'a'.repeat(64), secret)).to.equal(false);
    expect(verifyWebhookSignature(payload, 't=12345', secret)).to.equal(false);
  });

  it('rejects wrong secrets', () => {
    const sig = signWebhookPayload(payload, secret);
    expect(verifyWebhookSignature(payload, sig, 'wrong')).to.equal(false);
    const header = signWebhookHeader(payload, secret);
    expect(verifyWebhookSignature(payload, header, 'wrong')).to.equal(false);
  });
});
