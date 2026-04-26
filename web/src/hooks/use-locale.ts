'use client';

import { useCallback, useTransition } from 'react';
import { useLocale as useNextIntlLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import { locales, type Locale } from '@/i18n';
import { legacyLocaleCookieName, primaryLocaleCookieName } from '@/i18n/config';

/**
 * Hook for managing application locale
 * Handles locale switching via cookie-based persistence and router refresh.
 */
export function useLocale() {
  const locale = useNextIntlLocale() as Locale;
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  /**
   * Switches the application locale.
   * This implementation persists the locale in both the current and legacy
   * cookie names, then refreshes the App Router tree so server translations
   * are fetched again.
   */
  const setLocale = useCallback((newLocale: Locale) => {
    if (!locales.includes(newLocale) || newLocale === locale) return;

    startTransition(() => {
      const cookieOptions = 'path=/;max-age=31536000;samesite=lax';

      document.cookie = `${primaryLocaleCookieName}=${encodeURIComponent(newLocale)};${cookieOptions}`;
      document.cookie = `${legacyLocaleCookieName}=${encodeURIComponent(newLocale)};${cookieOptions}`;
      document.documentElement.lang = newLocale;
      router.refresh();
    });
  }, [locale, router]);

  return {
    locale,
    setLocale,
    isPending,
  };
}
