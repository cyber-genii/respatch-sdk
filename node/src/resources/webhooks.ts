import { RespatchClient } from '../client';
import { CreateWebhookRequest } from '../types';

export class Webhooks {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  async list(): Promise<any[]> {
    return this.client.get<any[]>('/api/v1/developers/webhooks');
  }

  async create(data: CreateWebhookRequest): Promise<any> {
    return this.client.post<any>('/api/v1/developers/webhooks', data);
  }

  async retrieve(id: string): Promise<any> {
    return this.client.get<any>(`/api/v1/developers/webhooks/${id}`);
  }

  async update(id: string, data: Partial<CreateWebhookRequest>): Promise<any> {
    return this.client.put<any>(`/api/v1/developers/webhooks/${id}`, data);
  }

  async delete(id: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/developers/webhooks/${id}`);
  }

  async getEvents(id: string): Promise<any[]> {
    return this.client.get<any[]>(`/api/v1/developers/webhooks/${id}/events`);
  }
}
