import { describe, expect, it } from 'vitest';
import { resolveRequestLocale } from '@/i18n/resolve-locale';

describe('resolveRequestLocale', () => {
  it('prefers NEXT_LOCALE when present', () => {
    expect(
      resolveRequestLocale({
        cookieLocales: {
          NEXT_LOCALE: 'zh-Hans',
          locale: 'en-US',
        },
        acceptLanguage: 'en-US,en;q=0.9',
      })
    ).toBe('zh-Hans');
  });

  it('falls back to the legacy locale cookie', () => {
    expect(
      resolveRequestLocale({
        cookieLocales: {
          locale: 'zh-Hans',
        },
        acceptLanguage: 'en-US,en;q=0.9',
      })
    ).toBe('zh-Hans');
  });

  it('maps short Accept-Language values to supported locales', () => {
    expect(
      resolveRequestLocale({
        acceptLanguage: 'zh-CN,zh;q=0.9,en;q=0.8',
      })
    ).toBe('zh-Hans');
  });

  it('falls back to the default locale for unsupported values', () => {
    expect(
      resolveRequestLocale({
        cookieLocales: {
          NEXT_LOCALE: 'fr-FR',
        },
        acceptLanguage: 'fr-FR,fr;q=0.9',
      })
    ).toBe('en-US');
  });
});
