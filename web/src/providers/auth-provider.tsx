'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/store/auth-store';

interface AuthProviderProps {
  children: React.ReactNode;
}

/**
 * Authentication provider that initializes auth state on app startup.
 * 
 * This provider initializes auth in the background without blocking page rendering.
 * Protected routes should use AuthGuard to enforce authentication.
 * Public routes (site, auth) will render immediately while auth state is being determined.
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const initializeAuth = useAuthStore.use.initializeAuth();

  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  // Render children immediately without waiting for auth initialization
  // Protected routes will use AuthGuard for blocking behavior
  return <>{children}</>;
}
