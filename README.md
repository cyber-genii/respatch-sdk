# Respatch SDKs

Official client libraries for the [Respatch](https://respatch.com) logistics API.

## Languages

| SDK | Path | Install |
|-----|------|---------|
| Node.js / TypeScript | `node/` | `npm install @respatch/node-sdk` |
| Python | `python/` | `pip install -e python/` (package `respatch`) |
| Go | `go/` | `go get github.com/transmit/transmit-go-sdk` (module path legacy; API is Respatch) |
| Rust | `rust/` | `cargo add respatch-rust-sdk --path rust/` |

> **Node rename:** `@transmit/node-sdk` is replaced by **`@respatch/node-sdk`**. The `Respatch` class and webhook helpers are unchanged; `Transmit` / `TransmitClient` remain deprecated aliases.

## Authentication

- **Integration calls** (delivery orders): `Authorization: Bearer <api_key>` (`rsp_*` keys from the developer dashboard).
- **Developer portal calls** (manage keys/webhooks/sandbox): JWT session from `/auth/login`.

Set `RESPATCH_API_KEY` in your environment (legacy `TRANSMIT_API_KEY` also supported in Node/Python).

## API response envelope

The backend wraps payloads as:

```json
{ "success": true, "message": "...", "data": <payload> }
```

All SDK HTTP clients unwrap `.data` automatically.

## Quick start (Node)

```typescript
import { Respatch } from '@respatch/node-sdk';

const respatch = new Respatch({
  apiKey: process.env.RESPATCH_API_KEY!,
  environment: 'sandbox', // sandbox-api.respatch.com — use 'production' for api.respatch.com
});

// Step 1: lock pricing
const quote = await respatch.orders.quote({
  pickup: { lat: 6.5244, lng: 3.3792 },
  dropoff: { address: '12 Admiralty Way, Lekki Phase 1, Lagos' },
  delivery_type: 'express',
  packages: [
    { weight_kg: 2, dimensions_cm: { length: 30, width: 20, height: 5 }, quantity: 1 },
  ],
});

// Step 2: create order (contacts + package metadata only)
const order = await respatch.orders.bookFromQuote({
  quoteId: quote.quote_id,
  pickup: { contact_name: 'Jane Sender', contact_phone: '+2348012345678' },
  dropoff: { contact_name: 'John Receiver', contact_phone: '+2348098765432' },
  packages: [{ description: 'Documents', value: 5000 }],
});

console.log(order.tracking_url, order.fare.breakdown);
```

`bookFromQuote` is a single HTTP call to Create Order (quote must already exist). Use `orders.quote()` first.

## Base URLs

| Frontier | Host | SDK `environment` |
|----------|------|-------------------|
| Production (`respatch-api`) | `https://api.respatch.com` | `production` (default) |
| Sandbox (`respatch-sandbox`) | `https://sandbox-api.respatch.com` | `sandbox` |

Override with `baseUrl` / `base_url` in client options. Do **not** point SDKs at the private platform backend — only `api.respatch.com` (respatch-api) and `sandbox-api.respatch.com` (respatch-sandbox).

## Canonical endpoints

Use **Delivery Orders** (`/api/v1/delivery-orders/quote` then `/api/v1/delivery-orders`). Legacy `/api/v1/api-deliveries` routes return migration errors.

## Webhook verification

Verify `X-Respatch-Signature` (HMAC-SHA256 hex of the raw body). Legacy `X-Transmit-Signature` is accepted by SDK helpers.

Webhook payloads for delivery status events include `package_id` for per-package tracking. **Partial fulfillment is supported** — one package may reach `exception` while others are `delivered`; rollup `order.status` becomes `partially_delivered`.

**Node**
```typescript
import { verifyWebhookSignature, getWebhookSignatureHeader } from '@respatch/node-sdk';
```

**Python**
```python
from respatch import verify_webhook_signature, get_webhook_signature_header
# legacy: from transmit import ...
```

**Go**
```go
import "github.com/transmit/transmit-go-sdk/transmit"
ok := transmit.VerifyWebhookSignature(body, sigHeader, secret)
```

**Rust**
```rust
use respatch_rust_sdk::verify_webhook_signature;
```
