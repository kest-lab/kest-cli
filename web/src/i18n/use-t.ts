import { locales, defaultLocale, type Locale } from './config';

// Basic translations for UI components
const translations: Record<Locale, any> = {
    'zh-Hans': {
        common: {
            datePlaceholder: '请选择日期',
            hour: '时',
            minute: '分',
            second: '秒',
            now: '现在',
            confirm: '确认',
            year: '年',
            month: '月'
        }
    },
    'en-US': {
        common: {
            datePlaceholder: 'Select date',
            hour: 'H',
            minute: 'M',
            second: 'S',
            now: 'Now',
            confirm: 'Confirm',
            year: 'Y',
            month: 'M'
        }
    }
};

/**
 * Basic translation hook shim.
 */
export function useT(scope?: string): ((key: string, variables?: Record<string, any>) => any) & Record<string, (key: string, variables?: Record<string, any>) => any> {
    // Simple locale detection from cookie
    const getLocale = (): Locale => {
        if (typeof document === 'undefined') return defaultLocale;
        const match = document.cookie.match(/locale=([^;]+)/);
        const locale = match?.[1] as Locale;
        return locales.includes(locale) ? locale : defaultLocale;
    };

    const locale = getLocale();
    const messages = translations[locale] || translations[defaultLocale];

    const t = (key: string, variables?: Record<string, any>) => {
        const keys = key.split('.');
        let value = scope ? (messages[scope] || {}) : messages;

        for (const k of keys) {
            if (value && typeof value === 'object') {
                value = value[k];
            } else {
                return key;
            }
        }

        if (typeof value === 'string' && variables) {
            Object.entries(variables).forEach(([k, v]) => {
                value = (value as string).replace(`{${k}}`, String(v));
            });
        }

        return value || key;
    };

    // Add common modules for compatibility with useT().common('key')
    const proxy = new Proxy(t, {
        get(target, prop) {
            if (typeof prop === 'string' && !['arguments', 'caller', 'prototype', 'name'].includes(prop)) {
                return (key: string, vars?: Record<string, any>) => t(`${prop}.${key}`, vars);
            }
            return (target as any)[prop];
        }
    });

    return proxy as any;
}
