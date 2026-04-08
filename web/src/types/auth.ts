export interface ApiUser {
  id: number;
  username: string;
  email: string;
  nickname?: string;
  avatar?: string;
  phone?: string;
  bio?: string;
  status: number | string;
  is_super_admin?: boolean;
  last_login?: string | null;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  user: ApiUser;
}

export interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  nickname?: string;
  phone?: string;
}

export interface PasswordResetRequest {
  email: string;
}

export interface PasswordResetResponse {
  message: string;
}

export interface UpdateProfileRequest {
  nickname?: string;
  avatar?: string;
  phone?: string;
  bio?: string;
}

export interface ChangePasswordRequest {
  old_password: string;
  new_password: string;
}

export interface MessageResponse {
  message: string;
}

export interface UserListParams {
  page?: number;
  perPage?: number;
}

export interface PaginationMeta {
  current_page: number;
  per_page: number;
  total: number;
  last_page: number;
  from: number;
  to: number;
}

export interface PaginationLinks {
  first: string;
  last: string;
  prev?: string | null;
  next?: string | null;
}

export interface PaginatedUsersResponse {
  items: ApiUser[];
  meta: PaginationMeta;
  links: PaginationLinks;
}

