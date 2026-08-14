import { RespatchClient } from '../client';

export class Vehicles {
  private client: RespatchClient;

  constructor(client: RespatchClient) {
    this.client = client;
  }

  async getTypes(): Promise<any[]> {
    return this.client.get<any[]>('/api/v1/vehicles/types');
  }
}
