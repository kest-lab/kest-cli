import request from '@/http';
import type {
  ApiUser,
  ChangePasswordRequest,
  LoginRequest,
  LoginResponse,
  MessageResponse,
  PaginatedUsersResponse,
  PasswordResetRequest,
  PasswordResetResponse,
  RegisterRequest,
  UpdateProfileRequest,
  UserListParams,
} from '@/types/auth';

interface RawPaginatedUsersResponse {
  code: number;
  message: string;
  data: ApiUser[];
  meta: PaginatedUsersResponse['meta'];
  links: PaginatedUsersResponse['links'];
}

export const authService = {
  login: (data: LoginRequest) => request.post<LoginResponse>('/login', data, { skipAuth: true }),

  register: (data: RegisterRequest) => request.post<ApiUser>('/register', data, { skipAuth: true }),

  requestPasswordReset: (data: PasswordResetRequest) =>
    request.post<PasswordResetResponse>('/password/reset', data, { skipAuth: true }),

  getProfile: () => request.get<ApiUser>('/users/profile'),

  updateProfile: (data: UpdateProfileRequest) => request.put<ApiUser>('/users/profile', data),

  changePassword: (data: ChangePasswordRequest) =>
    request.put<MessageResponse>('/users/password', data),

  deleteAccount: () => request.delete<void>('/users/account'),
};

export const userService = {
  list: async ({ page = 1, perPage = 10 }: UserListParams = {}): Promise<PaginatedUsersResponse> => {
    // 用户列表接口除了 data 之外，还会在响应顶层返回 meta / links，
    // 因此这里关闭默认 data 解包，保留完整响应结构再做一次前端归一化。
    const response = await request.get<RawPaginatedUsersResponse>('/users', {
      params: { page, per_page: perPage },
      unwrapData: false,
    });

    return {
      items: response.data,
      meta: response.meta,
      links: response.links,
    };
  },

  search: (query: string, limit = 20) =>
    // 后端实际提供的是搜索接口，作为列表页的补充能力。
    request.get<ApiUser[]>('/users/search', {
      params: { q: query, limit },
    }),

  getById: (id: number | string) => request.get<ApiUser>(`/users/${id}`),

  getInfo: (id: number | string) => request.get<ApiUser>(`/users/${id}/info`),
};

export type AuthService = typeof authService;
export type UserService = typeof userService;
