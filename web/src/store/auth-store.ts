import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import { createSelectors } from './utils/selectors';
import { authApi } from '@/services/auth';
import type { User, SystemFeatures, SetupStatus } from '@/types/auth';

interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    accessToken: string | null;
    systemFeatures: SystemFeatures | null;
    setupStatus: SetupStatus | null;
    isSystemReady: boolean;
    error: string | null;
    setUser: (user: User | null) => void;
    setLoading: (loading: boolean) => void;
    setError: (error: string | null) => void;
    setSystemFeatures: (features: SystemFeatures | null) => void;
    setSetupStatus: (status: SetupStatus | null) => void;
    refreshProfile: () => Promise<void>;
    initializeSystem: () => Promise<void>;
    clearAuth: () => void;
    reset: () => void;
}

const defaultState = {
    user: null,
    isAuthenticated: false,
    isLoading: false,
    accessToken: null,
    systemFeatures: null,
    setupStatus: null,
    isSystemReady: false,
    error: null,
};

const useAuthStoreBase = create<AuthState>()(
    persist(
        (set, get) => ({
            ...defaultState,
            setUser: (user) => set({
                user,
                isAuthenticated: !!user,
                error: null
            }),
            setLoading: (isLoading) => set({ isLoading }),
            setError: (error) => set({ error }),
            setSystemFeatures: (systemFeatures) => set({
                systemFeatures,
                isSystemReady: true
            }),
            setSetupStatus: (setupStatus) => set({ setupStatus }),
            refreshProfile: async () => {
                set({ isLoading: true });
                try {
                    const user = await authApi.getProfile();
                    set({ user, isLoading: false, isAuthenticated: true });
                } catch (error: any) {
                    set({ isLoading: false, error: error.message });
                }
            },
            initializeSystem: async () => {
                set({ isLoading: true });
                try {
                    // You can add more initialization logic here if needed
                    set({ isSystemReady: true });
                } catch (error) {
                    console.error('Failed to initialize system:', error);
                    set({ isSystemReady: true });
                } finally {
                    set({ isLoading: false });
                }
            },
            clearAuth: () => {
                set({
                    user: null,
                    isAuthenticated: false,
                    accessToken: null,
                    error: null,
                });
                localStorage.removeItem('auth-storage');
            },
            reset: () => set(defaultState),
        }),
        {
            name: 'auth-storage',
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({
                user: state.user,
                isAuthenticated: state.isAuthenticated,
                accessToken: state.accessToken,
                systemFeatures: state.systemFeatures,
                setupStatus: state.setupStatus,
                isSystemReady: state.isSystemReady,
            }),
        }
    )
);

export const getAuthTokens = () => {
    try {
        const state = useAuthStoreBase.getState();
        return { accessToken: state.accessToken };
    } catch {
        return { accessToken: null };
    }
};

export const setAuthTokens = (accessToken: string | null) => {
    useAuthStoreBase.setState({ accessToken, isAuthenticated: !!accessToken });
};

export const useAuthStore = createSelectors(useAuthStoreBase);

export const authSelectors = {
    isAdmin: (state: AuthState) => state.user?.role === 'admin',
    needsSetup: (state: AuthState) => state.setupStatus?.step === 'not_started',
    canRegister: (state: AuthState) => state.systemFeatures?.is_allow_register ?? true,
    hasEmailLogin: (state: AuthState) => state.systemFeatures?.enable_email_password_login ?? true,
    hasSocialLogin: (state: AuthState) => state.systemFeatures?.enable_social_oauth_login ?? false,
};

export const useCurrentUser = () => useAuthStore.use.user();
export const useIsAuthenticated = () => useAuthStore.use.isAuthenticated();
export const useAuthLoading = () => useAuthStore.use.isLoading();
export const useAuthError = () => useAuthStore.use.error();
