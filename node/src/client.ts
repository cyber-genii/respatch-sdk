import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';
import { unwrapData } from './response';
import { ClientOptions } from './types/index';

const PRODUCTION_URL = 'https://api.respatch.com';
const SANDBOX_URL = 'https://sandbox-api.respatch.com';

export class RespatchClient {
  private httpClient: AxiosInstance;

  constructor(options: ClientOptions) {
    const apiKey =
      options.apiKey ||
      process.env.RESPATCH_API_KEY ||
      process.env.TRANSMIT_API_KEY;

    if (!apiKey) {
      throw new Error('API Key is required (pass apiKey or set RESPATCH_API_KEY)');
    }

    const defaultBaseUrl =
      options.environment === 'sandbox' ? SANDBOX_URL : PRODUCTION_URL;
    const baseURL = options.baseUrl || defaultBaseUrl;

    this.httpClient = axios.create({
      baseURL,
      headers: {
        Authorization: `Bearer ${apiKey}`,
        'Content-Type': 'application/json',
        'User-Agent': 'Respatch-Node-SDK/1.1.0',
      },
      timeout: 30000,
    });

    this.httpClient.interceptors.response.use(
      (response: AxiosResponse) => response,
      (error: AxiosError) => {
        let errorMessage = 'An unknown error occurred during the API request';
        if (error.response && error.response.data) {
          const data = error.response.data as { message?: string; error?: string };
          errorMessage = data.message || data.error || errorMessage;
        } else if (error.message) {
          errorMessage = error.message;
        }
        throw new Error(
          `Respatch API Error [${error.response?.status || 'Unknown Status'}]: ${errorMessage}`,
        );
      },
    );
  }

  async get<T>(path: string, params?: Record<string, unknown>): Promise<T> {
    const response = await this.httpClient.get(path, { params });
    return unwrapData<T>(response.data);
  }

  async post<T>(path: string, data?: unknown): Promise<T> {
    const response = await this.httpClient.post(path, data);
    return unwrapData<T>(response.data);
  }

  async put<T>(path: string, data?: unknown): Promise<T> {
    const response = await this.httpClient.put(path, data);
    return unwrapData<T>(response.data);
  }

  async delete<T>(path: string): Promise<T> {
    const response = await this.httpClient.delete(path);
    return unwrapData<T>(response.data);
  }
}

export { PRODUCTION_URL as RESPATCH_PRODUCTION_URL, SANDBOX_URL as RESPATCH_SANDBOX_URL };

/** @deprecated Use `RespatchClient` */
export { RespatchClient as TransmitClient };
