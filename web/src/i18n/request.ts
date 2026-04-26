import { getRequestConfig } from 'next-intl/server';
import { cookies, headers } from 'next/headers';
import { legacyLocaleCookieName, primaryLocaleCookieName } from './config';
import { loadAllModules } from './loader';
import { resolveRequestLocale } from './resolve-locale';

export default getRequestConfig(async () => {
  const cookieStore = await cookies();
  const headerStore = await headers();
  const locale = resolveRequestLocale({
    cookieLocales: {
      [primaryLocaleCookieName]: cookieStore.get(primaryLocaleCookieName)?.value,
      [legacyLocaleCookieName]: cookieStore.get(legacyLocaleCookieName)?.value,
    },
    acceptLanguage: headerStore.get('accept-language'),
  });

  return {
    locale,
    messages: await loadAllModules(locale),
  };
});
