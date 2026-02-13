import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from 'axios';
import { handleError } from './error-handler';
import { getAuthTokens } from '@/store/auth-store';
import { env } from '@/config/env';

export interface RequestConfig extends AxiosRequestConfig {
    skipAuth?: boolean;
    skipErrorHandler?: boolean;
}

export class ApiError extends Error {
    code: string | number;
    status?: number;
    constructor(message: string, code: string | number, status?: number) {
        super(message);
        this.code = code;
        this.status = status;
    }
}

class HttpClient {
    private instance: AxiosInstance;

    constructor(config: RequestConfig) {
        this.instance = axios.create({
            timeout: 30000,
            ...config,
        });

        this.setupInterceptors();
    }

    private setupInterceptors() {
        this.instance.interceptors.request.use(
            (config) => {
                const { skipAuth } = config as RequestConfig;
                if (!skipAuth) {
                    const { accessToken } = getAuthTokens();
                    if (accessToken) {
                        config.headers.Authorization = "Bearer " + accessToken;
                    }
                }
                return config;
            },
            (error) => Promise.reject(error)
        );

        this.instance.interceptors.response.use(
            (response) => {
                const res = response.data;
                // Standard API response unwrap
                if (res && typeof res === 'object' && 'code' in res) {
                    if (res.code === 0) {
                        return res.data;
                    }
                    // Business logic error with HTTP 200
                    const apiError = new ApiError(
                        res.message || 'Unknown Error',
                        res.code,
                        response.status
                    );
                    handleError(apiError);
                    return Promise.reject(apiError);
                }
                return res;
            },
            (error: AxiosError) => {
                const { skipErrorHandler } = error.config as RequestConfig;
                const body = error.response?.data as any;

                const apiError = new ApiError(
                    body?.message || body?.error || error.message,
                    body?.code || 'FETCH_ERROR',
                    error.response?.status
                );

                if (!skipErrorHandler) {
                    handleError(apiError);
                }

                return Promise.reject(apiError);
            }
        );
    }

    public get<T = any>(url: string, config?: RequestConfig): Promise<T> {
        return this.instance.get(url, config);
    }

    public post<T = any>(url: string, data?: any, config?: RequestConfig): Promise<T> {
        return this.instance.post(url, data, config);
    }

    public put<T = any>(url: string, data?: any, config?: RequestConfig): Promise<T> {
        return this.instance.put(url, data, config);
    }

    public patch<T = any>(url: string, data?: any, config?: RequestConfig): Promise<T> {
        return this.instance.patch(url, data, config);
    }

    public delete<T = any>(url: string, config?: RequestConfig): Promise<T> {
        return this.instance.delete(url, config);
    }
}

export const createRequest = (config: RequestConfig = {}) => {
    return new HttpClient(config);
};

export const request = createRequest({
    baseURL: env.VITE_API_URL,
});

export default request;
