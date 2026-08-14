import { RespatchClient } from '../client';
import { unwrapList } from '../response';
import {
  CalculateOrderFareRequest,
  CreateDeliveryOrderRequest,
  CreateDeliveryOrderResponse,
} from '../types/deliveryOrders';

/**
 * @deprecated Use `Orders` (`client.orders`) — `/api/v1/api-deliveries` returns a Deprecation header.
 */
export class Deliveries {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  async create(data: CreateDeliveryOrderRequest): Promise<CreateDeliveryOrderResponse> {
    return this.client.post<CreateDeliveryOrderResponse>('/api/v1/delivery-orders', data);
  }

  async list(params?: Record<string, unknown>): Promise<CreateDeliveryOrderResponse[]> {
    return unwrapList<CreateDeliveryOrderResponse>(
      await this.client.get<unknown>('/api/v1/delivery-orders', params),
    );
  }

  async retrieve(deliveryId: string): Promise<CreateDeliveryOrderResponse> {
    return this.client.get<CreateDeliveryOrderResponse>(`/api/v1/delivery-orders/${deliveryId}`);
  }

  async quote(data: CalculateOrderFareRequest): Promise<unknown> {
    return this.client.post('/api/v1/delivery-orders/calculate-fare', data);
  }
}
