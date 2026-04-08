'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { authService } from '@/services/auth';
import { useAuthStore } from '@/store/auth-store';
import type {
  ChangePasswordRequest,
  LoginRequest,
  LoginResponse,
  PasswordResetRequest,
  RegisterRequest,
  UpdateProfileRequest,
} from '@/types/auth';

export const authKeys = {
  all: ['auth'] as const,
  profile: () => [...authKeys.all, 'profile'] as const,
};

export function useProfile(enabled = true) {
  return useQuery({
    queryKey: authKeys.profile(),
    queryFn: authService.getProfile,
    enabled,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  const setSession = useAuthStore.use.setSession();

  return useMutation({
    mutationFn: (data: LoginRequest) => authService.login(data),
    onSuccess: (result: LoginResponse) => {
      // 登录成功后同时更新本地会话和 profile 缓存，
      // 这样控制台和站点头部都能立即反映最新登录状态。
      setSession(result.user, result.access_token);
      queryClient.setQueryData(authKeys.profile(), result.user);
    },
  });
}

export function useRegister() {
  return useMutation({
    mutationFn: (data: RegisterRequest) => authService.register(data),
  });
}

export function useRequestPasswordReset() {
  return useMutation({
    mutationFn: (data: PasswordResetRequest) => authService.requestPasswordReset(data),
  });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  const updateUser = useAuthStore.use.updateUser();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => authService.updateProfile(data),
    onSuccess: (user) => {
      // 设置页更新成功后，同步刷新 store 和 query cache，
      // 避免头像、昵称等信息在不同区域出现短暂不一致。
      updateUser(user);
      queryClient.setQueryData(authKeys.profile(), user);
    },
  });
}

export function useChangePassword() {
  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => authService.changePassword(data),
  });
}

export function useDeleteAccount() {
  const queryClient = useQueryClient();
  const clearSession = useAuthStore.use.clearSession();

  return useMutation({
    mutationFn: authService.deleteAccount,
    onSuccess: async () => {
      // 账号删除后立即清空本地登录态，防止页面继续携带失效 token。
      clearSession();
      await queryClient.invalidateQueries({ queryKey: authKeys.all });
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const clearSession = useAuthStore.use.clearSession();

  return () => {
    clearSession();
    queryClient.removeQueries({ queryKey: authKeys.all });
  };
}
