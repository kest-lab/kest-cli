export interface RequestKeyValue {
  key: string;
  value: string;
  type?: string;
  enabled?: boolean;
  description?: string;
}

export interface RequestBasicAuth {
  username: string;
  password: string;
}

export interface RequestBearerAuth {
  token: string;
}

export interface RequestApiKeyAuth {
  key: string;
  value: string;
  in?: string;
  add_to?: string;
}

export interface RequestOAuth2Auth {
  grant_type: string;
  auth_url?: string;
  token_url?: string;
  client_id?: string;
  client_secret?: string;
  scope?: string;
  username?: string;
  password?: string;
}

export interface RequestAuthConfig {
  type: string;
  basic?: RequestBasicAuth;
  bearer?: RequestBearerAuth;
  api_key?: RequestApiKeyAuth;
  oauth2?: RequestOAuth2Auth;
}

export interface ProjectRequest {
  id: number;
  collection_id: number;
  name: string;
  description: string;
  method: string;
  url: string;
  headers: RequestKeyValue[];
  query_params: RequestKeyValue[];
  path_params: Record<string, string>;
  body: string;
  body_type: string;
  auth?: RequestAuthConfig | null;
  pre_request?: string;
  test?: string;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface RequestListMeta {
  total: number;
  page: number;
  per_page: number;
  pages: number;
}

export interface RequestListResponse {
  items: ProjectRequest[];
  meta: RequestListMeta;
}

export interface CreateRequestRequest {
  collection_id: number;
  name: string;
  description?: string;
  method: string;
  url: string;
  headers?: RequestKeyValue[];
  query_params?: RequestKeyValue[];
  path_params?: Record<string, string>;
  body?: string;
  body_type?: string;
  auth?: RequestAuthConfig | null;
  pre_request?: string;
  test?: string;
  sort_order?: number;
}

export interface UpdateRequestRequest {
  name?: string;
  description?: string;
  method?: string;
  url?: string;
  headers?: RequestKeyValue[];
  query_params?: RequestKeyValue[];
  path_params?: Record<string, string>;
  body?: string;
  body_type?: string;
  auth?: RequestAuthConfig | null;
  pre_request?: string;
  test?: string;
  sort_order?: number;
}
