import { env } from './env';

const normalizeBasePath = (value: string) => {
  if (!value) {
    return '';
  }

  const withLeadingSlash = value.startsWith('/') ? value : `/${value}`;
  return withLeadingSlash.replace(/\/+$/, '');
};

const apiOrigin = env.NEXT_PUBLIC_API_URL.replace(/\/$/, '');

// 将接口版本前缀单独抽出来，后续切换 /v2 或网关前缀时只需要改 env。
export const apiBasePath = normalizeBasePath(env.NEXT_PUBLIC_API_BASE_PATH);
export const apiProxyPath = normalizeBasePath(env.NEXT_PUBLIC_API_PROXY_PATH);
export const apiUseProxy = env.NEXT_PUBLIC_API_USE_PROXY;

// 外部真实 API 地址，用于服务端转发和调试展示。
export const apiExternalBaseUrl = apiOrigin ? `${apiOrigin}${apiBasePath}` : apiBasePath;

// 浏览器默认通过 Next 同源代理访问，绕开线上 API 当前的 CORS 限制。
export const apiBaseUrl = apiUseProxy ? apiProxyPath : apiExternalBaseUrl;

export const buildApiPath = (path: string) => {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${apiBasePath}${normalizedPath}`;
};

export const buildApiUrl = (path: string) => {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  if (apiUseProxy) {
    return `${apiProxyPath}${normalizedPath}`;
  }

  const apiPath = buildApiPath(normalizedPath);
  return apiOrigin ? `${apiOrigin}${apiPath}` : apiPath;
};
