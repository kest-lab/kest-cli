import {
  defaultLocale,
  legacyLocaleCookieName,
  localeMapping,
  locales,
  primaryLocaleCookieName,
  type Locale,
} from './config';

interface ResolveLocaleInput {
  cookieLocales?: Partial<Record<typeof primaryLocaleCookieName | typeof legacyLocaleCookieName, string | undefined>>;
  acceptLanguage?: string | null;
}

function normalizeLocale(candidate?: string): Locale | undefined {
  if (!candidate) {
    return undefined;
  }

  try {
    const decoded = decodeURIComponent(candidate).trim();

    if (locales.includes(decoded as Locale)) {
      return decoded as Locale;
    }

    if (localeMapping[decoded]) {
      return localeMapping[decoded];
    }

    const languageOnly = decoded.split('-')[0];
    if (localeMapping[languageOnly]) {
      return localeMapping[languageOnly];
    }
  } catch {
    return undefined;
  }

  return undefined;
}

export function resolveRequestLocale({
  cookieLocales,
  acceptLanguage,
}: ResolveLocaleInput): Locale {
  const cookieCandidates = [
    cookieLocales?.[primaryLocaleCookieName],
    cookieLocales?.[legacyLocaleCookieName],
  ];

  for (const candidate of cookieCandidates) {
    const locale = normalizeLocale(candidate);
    if (locale) {
      return locale;
    }
  }

  if (acceptLanguage) {
    const languages = acceptLanguage.split(',').map((part) => part.split(';')[0]?.trim());

    for (const candidate of languages) {
      const locale = normalizeLocale(candidate);
      if (locale) {
        return locale;
      }
    }
  }

  return defaultLocale;
}
