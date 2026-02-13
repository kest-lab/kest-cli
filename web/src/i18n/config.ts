export const locales = ['zh-Hans', 'en-US'] as const;
export type Locale = (typeof locales)[number];
export const defaultLocale: Locale = 'zh-Hans';

export const localeNames: Record<Locale, string> = {
    'zh-Hans': '简体中文',
    'en-US': 'English',
};

export const localeMapping: Record<string, Locale> = {
    'zh': 'zh-Hans',
    'zh-CN': 'zh-Hans',
    'en': 'en-US',
    'en-US': 'en-US',
};

export const isLocaleSwitcherEnabled = true;
