// Supported locales
export const locales = ['zh-Hans', 'en-US'] as const;
export type Locale = (typeof locales)[number];

// Default locale
export const defaultLocale: Locale = 'en-US';

// Switches
export const isLocaleSwitcherEnabled = true;

export const localeNames: Record<Locale, string> = {
  'zh-Hans': '简体中文',
  'en-US': 'English',
};

// Mapping for Accept-Language header
export const localeMapping: Record<string, Locale> = {
  'zh': 'zh-Hans',
  'zh-CN': 'zh-Hans',
  'en': 'en-US',
  'en-US': 'en-US',
};
