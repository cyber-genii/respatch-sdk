import { RespatchClient } from './client';
import { ApiKeys } from './resources/apiKeys';
import { Deliveries } from './resources/deliveries';
import { Orders } from './resources/orders';
import { Sandbox } from './resources/sandbox';
import { Vehicles } from './resources/vehicles';
import { Webhooks } from './resources/webhooks';
import { ClientOptions } from './types';

export class Respatch {
  public orders: Orders;
  /** @deprecated Use `orders` — kept for backward compatibility. */
  public deliveries: Deliveries;
  public apiKeys: ApiKeys;
  public webhooks: Webhooks;
  public sandbox: Sandbox;
  public vehicles: Vehicles;
  private client: RespatchClient;

  constructor(options: ClientOptions) {
    this.client = new RespatchClient(options);
    this.orders = new Orders(this.client);
    this.deliveries = new Deliveries(this.client);
    this.apiKeys = new ApiKeys(this.client);
    this.webhooks = new Webhooks(this.client);
    this.sandbox = new Sandbox(this.client);
    this.vehicles = new Vehicles(this.client);
  }
}

/** @deprecated Use `Respatch` */
export { Respatch as Transmit };

export * from './types/deliveryOrders';
export {
  verifyWebhookSignature,
  signWebhookPayload,
  getWebhookSignatureHeader,
  RESPATCH_SIGNATURE_HEADER,
} from './webhook';
export {
  RespatchClient,
  TransmitClient,
  RESPATCH_PRODUCTION_URL,
  RESPATCH_SANDBOX_URL,
} from './client';
