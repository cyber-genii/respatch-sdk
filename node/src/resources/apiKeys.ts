import { RespatchClient } from '../client';
import { unwrapList, unwrapPaginated } from '../response';
import { CreateApiKeyRequest } from '../types';

export class ApiKeys {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  async list(): Promise<any[]> {
    const raw = await this.client.get<unknown>('/api/v1/developers/api-keys');
    return unwrapPaginated(raw).items;
  }

  async create(data: CreateApiKeyRequest): Promise<any> {
    return this.client.post<any>('/api/v1/developers/api-keys', data);
  }

  async retrieve(id: string): Promise<any> {
    return this.client.get<any>(`/api/v1/developers/api-keys/${id}`);
  }

  async revoke(id: string): Promise<void> {
    return this.client.post<void>(`/api/v1/developers/api-keys/${id}/revoke`);
  }

  async delete(id: string): Promise<void> {
    return this.client.delete<void>(`/api/v1/developers/api-keys/${id}`);
  }

  async getUsage(id: string): Promise<any> {
    return this.client.get<any>(`/api/v1/developers/api-keys/${id}/usage`);
  }
}
