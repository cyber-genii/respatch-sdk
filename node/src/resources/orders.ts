import { RespatchClient } from '../client';
import { unwrapList } from '../response';
import {
  BookFromQuoteParams,
  CreateDeliveryOrderRequest,
  CreateDeliveryOrderResponse,
  DeliveryOrderTrackResponse,
  GetQuoteRequest,
  GetQuoteResponse,
} from '../types/deliveryOrders';

export class Orders {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  /** Get a persisted fare quote (required before creating an order). */
  async quote(data: GetQuoteRequest): Promise<GetQuoteResponse> {
    return this.client.post<GetQuoteResponse>('/api/v1/delivery-orders/quote', data);
  }

  /** @deprecated Use `quote()` — hits deprecated `/calculate-fare` alias. */
  async calculateFare(data: GetQuoteRequest): Promise<GetQuoteResponse> {
    return this.client.post<GetQuoteResponse>('/api/v1/delivery-orders/calculate-fare', data);
  }

  async create(data: CreateDeliveryOrderRequest): Promise<CreateDeliveryOrderResponse> {
    return this.client.post<CreateDeliveryOrderResponse>('/api/v1/delivery-orders', data);
  }

  /**
   * Create an order from a prior quote (single HTTP call to Create Order).
   * Obtain `quoteId` first via `quote()`.
   */
  async bookFromQuote(params: BookFromQuoteParams): Promise<CreateDeliveryOrderResponse> {
    const { quoteId, pickup, dropoff, packages, ...rest } = params;
    return this.create({
      quote_id: quoteId,
      pickup,
      dropoff,
      packages,
      ...rest,
    });
  }

  async list(params?: Record<string, unknown>): Promise<CreateDeliveryOrderResponse[]> {
    return unwrapList<CreateDeliveryOrderResponse>(
      await this.client.get<unknown>('/api/v1/delivery-orders', params),
    );
  }

  async retrieve(orderId: string): Promise<CreateDeliveryOrderResponse> {
    return this.client.get<CreateDeliveryOrderResponse>(`/api/v1/delivery-orders/${orderId}`);
  }

  async track(orderId: string): Promise<DeliveryOrderTrackResponse> {
    return this.client.get<DeliveryOrderTrackResponse>(`/api/v1/delivery-orders/${orderId}/track`);
  }

  async cancel(orderId: string): Promise<void> {
    await this.client.post(`/api/v1/delivery-orders/${orderId}/cancel`);
  }
}
