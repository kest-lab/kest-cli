import { getRequestConfig } from 'next-intl/server';
import { cookies, headers } from 'next/headers';
import { defaultLocale, locales, localeMapping, type Locale } from './config';
import { loadAllModules } from './loader';

export default getRequestConfig(async () => {
  // Try to get locale from cookie first
  const cookieStore = await cookies();
  const localeCookie = cookieStore.get('locale')?.value as Locale | undefined;

  // Determine locale with fallback logic
  let locale: Locale = defaultLocale;

  if (localeCookie && locales.includes(localeCookie)) {
    locale = localeCookie;
  } else {
    // Try Accept-Language header
    const headerStore = await headers();
    const acceptLanguage = headerStore.get('accept-language');
    
    if (acceptLanguage) {
      // Parse Accept-Language header (e.g., "zh-CN,zh;q=0.9,en;q=0.8")
      const languages = acceptLanguage.split(',').map(lang => lang.split(';')[0].trim());
      
      for (const lang of languages) {
        // Try exact match first, then language-only match
        if (localeMapping[lang]) {
          locale = localeMapping[lang];
          break;
        }
        const langOnly = lang.split('-')[0];
        if (localeMapping[langOnly]) {
          locale = localeMapping[langOnly];
          break;
        }
      }
    }
  }

  return {
    locale,
    messages: await loadAllModules(locale),
  };
});
