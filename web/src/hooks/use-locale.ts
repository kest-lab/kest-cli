import { useCallback, useTransition } from 'react';
import { locales, type Locale } from '@/i18n';

/**
 * Hook for managing application locale.
 * Adapted for Vite/React without next-intl.
 */
export function useLocale() {
    // We'll need a way to get the current locale.
    // For now, we'll read it from a cookie or default to the first locale.
    const getLocaleFromCookie = () => {
        const match = document.cookie.match(/locale=([^;]+)/);
        return (match?.[1] as Locale) || locales[0];
    };

    const locale = getLocaleFromCookie();
    const [isPending, startTransition] = useTransition();

    const setLocale = useCallback((newLocale: Locale) => {
        if (!locales.includes(newLocale)) return;

        startTransition(() => {
            // Set cookie and reload to apply new locale
            document.cookie = `locale=${newLocale};path=/;max-age=31536000`;
            window.location.reload();
        });
    }, []);

    return {
        locale,
        setLocale,
        isPending,
    };
}
