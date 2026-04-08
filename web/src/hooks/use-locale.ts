'use client';

import { useCallback, useTransition } from 'react';
import { useLocale as useNextIntlLocale } from 'next-intl';
import { locales, type Locale } from '@/i18n';

/**
 * Hook for managing application locale
 * Handles locale switching via cookie-based persistence and page reload.
 */
export function useLocale() {
  const locale = useNextIntlLocale() as Locale;
  const [isPending, startTransition] = useTransition();

  /**
   * Switches the application locale.
   * This implementation uses a cookie to persist the choice and reloads the page
   * to ensure all server-side and client-side parts are updated.
   */
  const setLocale = useCallback((newLocale: Locale) => {
    if (!locales.includes(newLocale)) return;

    startTransition(() => {
      // Set cookie and reload to apply new locale
      // Note: 'locale' is the standard cookie name expected by next-intl middleware in this scaffold.
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
