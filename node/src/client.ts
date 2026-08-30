import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { unwrapData } from './response';
import { ClientOptions } from './types/index';

/** Public frontiers only — never the private platform (Koyeb) backend. */
const PRODUCTION_URL = 'https://api.respatch.com';
const SANDBOX_URL = 'https://sandbox-api.respatch.com';

const RETRYABLE_STATUSES = new Set([429, 503]);
const MAX_RETRIES = 4;
const INITIAL_BACKOFF_MS = 500;

export interface RequestOptions {
  idempotencyKey?: string;
}

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

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

  private async requestWithRetry<T>(
    config: AxiosRequestConfig,
  ): Promise<AxiosResponse<T>> {
    let backoff = INITIAL_BACKOFF_MS;
    for (let attempt = 0; attempt <= MAX_RETRIES; attempt++) {
      try {
        return await this.httpClient.request<T>(config);
      } catch (error) {
        const axiosError = error as AxiosError;
        const status = axiosError.response?.status;
        const retryable =
          status !== undefined && RETRYABLE_STATUSES.has(status);
        if (!retryable || attempt >= MAX_RETRIES) {
          throw error;
        }
        await sleep(backoff);
        backoff *= 2;
      }
    }
    throw new Error('Request failed after retries');
  }

  private mergeHeaders(options?: RequestOptions): Record<string, string> | undefined {
    if (!options?.idempotencyKey) {
      return undefined;
    }
    return { 'Idempotency-Key': options.idempotencyKey };
  }

  async get<T>(path: string, params?: Record<string, unknown>): Promise<T> {
    const response = await this.requestWithRetry<T>({ method: 'get', url: path, params });
    return unwrapData<T>(response.data);
  }

  async post<T>(
    path: string,
    data?: unknown,
    options?: RequestOptions,
  ): Promise<T> {
    const response = await this.requestWithRetry<T>({
      method: 'post',
      url: path,
      data,
      headers: this.mergeHeaders(options),
    });
    return unwrapData<T>(response.data);
  }

  async put<T>(path: string, data?: unknown): Promise<T> {
    const response = await this.requestWithRetry<T>({ method: 'put', url: path, data });
    return unwrapData<T>(response.data);
  }

  async delete<T>(path: string): Promise<T> {
    const response = await this.requestWithRetry<T>({ method: 'delete', url: path });
    return unwrapData<T>(response.data);
  }
}

export { PRODUCTION_URL as RESPATCH_PRODUCTION_URL, SANDBOX_URL as RESPATCH_SANDBOX_URL };

/** @deprecated Use `RespatchClient` */
export { RespatchClient as TransmitClient };
