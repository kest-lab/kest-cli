/**
 * String utility functions
 * Common operations for string manipulation
 */

/**
 * Truncates a string to a specified length and adds ellipsis if needed
 * @param str - String to truncate
 * @param length - Maximum length before truncation 
 * @param ellipsis - Optional custom ellipsis string (default '...')
 */
export function truncate(str: string, length: number, ellipsis = '...'): string {
    if (!str) return '';
    if (str.length <= length) return str;
    return str.substring(0, length).trim() + ellipsis;
}

/**
 * Capitalizes the first letter of a string
 * @param str - String to capitalize
 */
export function capitalize(str: string): string {
    if (!str) return '';
    return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Formats a string to title case (capitalizes the first letter of each word)
 * @param str - String to format
 */
export function toTitleCase(str: string): string {
    if (!str) return '';
    return str
        .toLowerCase()
        .split(' ')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
}

/**
 * Slugifies a string (converts to lowercase, replaces spaces with hyphens, removes special chars)
 * @param str - String to slugify
 */
export function slugify(str: string): string {
    if (!str) return '';
    return str
        .toLowerCase()
        .replace(/\s+/g, '-')
        .replace(/[^\w-]+/g, '')
        .replace(/--+/g, '-')
        .replace(/^-+/, '')
        .replace(/-+$/, '');
}

/**
 * Converts a string to camelCase
 * @param str - String to convert
 */
export function toCamelCase(str: string): string {
    if (!str) return '';
    return str
        .replace(/(?:^\w|[A-Z]|\b\w)/g, (letter, index) =>
            index === 0 ? letter.toLowerCase() : letter.toUpperCase()
        )
        .replace(/\s+/g, '');
}

/**
 * Format a number as currency
 * @param value - Number to format
 * @param options - Intl.NumberFormatOptions or locale string
 * @param currency - Currency code (default 'USD', only used if options is string)
 */
export function formatCurrency(
    value: number,
    options: Intl.NumberFormatOptions | string = 'en-US',
    currency = 'USD'
): string {
    if (typeof options === 'string') {
        return new Intl.NumberFormat(options, {
            style: 'currency',
            currency,
        }).format(value);
    }
    return new Intl.NumberFormat(undefined, {
        style: 'currency',
        currency,
        ...options
    }).format(value);
}

/**
 * Format a number as percentage
 * @param value - Number to format (0-1)
 * @param options - Intl.NumberFormatOptions or locale string
 * @param decimals - Number of decimal places (default 2, only used if options is string)
 */
export function formatPercent(
    value: number,
    options: Intl.NumberFormatOptions | string = 'en-US',
    decimals = 2
): string {
    if (typeof options === 'string') {
        return new Intl.NumberFormat(options, {
            style: 'percent',
            minimumFractionDigits: decimals,
            maximumFractionDigits: decimals,
        }).format(value);
    }
    return new Intl.NumberFormat(undefined, {
        style: 'percent',
        ...options
    }).format(value);
}
