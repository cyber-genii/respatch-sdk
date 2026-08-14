import { RespatchClient } from '../client';
import { CreateSandboxDataRequest } from '../types';

export class Sandbox {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  async createData(data: CreateSandboxDataRequest): Promise<any> {
    return this.client.post<any>('/api/v1/developers/sandbox', data);
  }

  async listData(): Promise<any[]> {
    return this.client.get<any[]>('/api/v1/developers/sandbox');
  }

  async clearData(): Promise<void> {
    return this.client.delete<void>('/api/v1/developers/sandbox');
  }
}
