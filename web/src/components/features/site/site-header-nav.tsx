/**
 * @component SiteHeaderNav
 */
'use client';

import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useAuthStore } from '@/store/auth-store';

export function SiteHeaderNav() {
    const isAuthenticated = useAuthStore.use.isAuthenticated();
    const user = useAuthStore.use.user();
    const isSystemReady = useAuthStore.use.isSystemReady();

    if (!isSystemReady) {
        return (
            <div className="flex items-center gap-4">
                <Link to="/login" title="Sign In">
                    <Button variant="ghost" size="sm">Sign In</Button>
                </Link>
                <Link to="/register" title="Get Started">
                    <Button size="sm">Get Started</Button>
                </Link>
            </div>
        );
    }

    if (isAuthenticated && user) {
        return (
            <div className="flex items-center gap-4">
                <Link to="/console" title="Dashboard">
                    <Button variant="ghost" size="sm">Dashboard</Button>
                </Link>
                <Link to="/console/settings" title={user.name || user.username || 'Profile'}>
                    <Button size="sm" variant="outline">{user.name || user.username || 'Profile'}</Button>
                </Link>
            </div>
        );
    }

    return (
        <div className="flex items-center gap-4">
            <Link to="/login" title="Sign In">
                <Button variant="ghost" size="sm">Sign In</Button>
            </Link>
            <Link to="/register" title="Get Started">
                <Button size="sm">Get Started</Button>
            </Link>
        </div>
    );
}
