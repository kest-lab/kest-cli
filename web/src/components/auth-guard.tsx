'use client';

import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/store/auth-store';
import { authConfig } from '@/config/auth';

interface AuthGuardProps {
    children: React.ReactNode;
    redirectTo?: string;
    showLoading?: boolean;
}

export function AuthGuard({
    children,
    redirectTo = authConfig.routes.login,
    showLoading = true,
}: AuthGuardProps) {
    const navigate = useNavigate();
    const location = useLocation();

    const isAuthenticated = useAuthStore.use.isAuthenticated();
    const isLoading = useAuthStore.use.isLoading();
    const isSystemReady = useAuthStore.use.isSystemReady();

    useEffect(() => {
        if (!isSystemReady) return;

        if (!isLoading && !isAuthenticated) {
            const returnUrl = encodeURIComponent(location.pathname + location.search);
            navigate(`${redirectTo}?returnUrl=${returnUrl}`, { replace: true });
        }
    }, [isAuthenticated, isLoading, isSystemReady, location, redirectTo, navigate]);

    if (!isSystemReady || isLoading) {
        if (showLoading) {
            return (
                <div className="flex min-h-screen items-center justify-center">
                    <div className="flex flex-col items-center space-y-4">
                        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
                        <p className="text-sm text-muted-foreground">Loading...</p>
                    </div>
                </div>
            );
        }
        return null;
    }

    if (!isAuthenticated) {
        return null;
    }

    return <>{children}</>;
}

export default AuthGuard;
