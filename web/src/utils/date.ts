/**
 * Date utility functions
 * Common operations for date handling and formatting
 */

/**
 * Format date according to specified format or locale options
 * @param date - Date to format
 * @param options - Intl.DateTimeFormatOptions or a format string (limited support for YYYY-MM-DD)
 * @param locale - Locale string (default: from browser/system)
 */
export function formatDate(
    date: Date | string | number,
    options: Intl.DateTimeFormatOptions | string = 'YYYY-MM-DD',
    locale?: string
): string {
    const d = new Date(date);

    if (isNaN(d.getTime())) {
        return 'Invalid Date';
    }

    if (typeof options === 'string') {
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');
        const hours = String(d.getHours()).padStart(2, '0');
        const minutes = String(d.getMinutes()).padStart(2, '0');
        const seconds = String(d.getSeconds()).padStart(2, '0');

        return options
            .replace('YYYY', String(year))
            .replace('MM', month)
            .replace('DD', day)
            .replace('HH', hours)
            .replace('mm', minutes)
            .replace('ss', seconds);
    }

    return new Intl.DateTimeFormat(locale, options).format(d);
}

/**
 * Get relative time string (e.g., "2 hours ago", "in 3 days")
 * @param date - Date to compare
 * @param baseDate - Base date to compare against (default: now)
 * @param locale - Locale string (default: 'en')
 */
export function getRelativeTime(
    date: Date | string | number,
    baseDate: Date | string | number = new Date(),
    locale = 'en'
): string {
    const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
    const d1 = new Date(date);
    const d2 = new Date(baseDate);

    const diffInMs = d1.getTime() - d2.getTime();
    const diffInSecs = Math.round(diffInMs / 1000);
    const diffInMins = Math.round(diffInSecs / 60);
    const diffInHours = Math.round(diffInMins / 60);
    const diffInDays = Math.round(diffInHours / 24);
    const diffInMonths = Math.round(diffInDays / 30);
    const diffInYears = Math.round(diffInDays / 365);

    if (Math.abs(diffInSecs) < 60) {
        return rtf.format(diffInSecs, 'second');
    } else if (Math.abs(diffInMins) < 60) {
        return rtf.format(diffInMins, 'minute');
    } else if (Math.abs(diffInHours) < 24) {
        return rtf.format(diffInHours, 'hour');
    } else if (Math.abs(diffInDays) < 30) {
        return rtf.format(diffInDays, 'day');
    } else if (Math.abs(diffInMonths) < 12) {
        return rtf.format(diffInMonths, 'month');
    } else {
        return rtf.format(diffInYears, 'year');
    }
}

/**
 * Check if a date is today
 * @param date - Date to check
 */
export function isToday(date: Date | string | number): boolean {
    const d = new Date(date);
    const today = new Date();
    return d.setHours(0, 0, 0, 0) === today.setHours(0, 0, 0, 0);
}

/**
 * Add time to a date
 * @param date - Base date
 * @param amount - Amount to add
 * @param unit - Unit to add (day, month, year, hour, minute, second)
 */
export function addTime(
    date: Date | string | number,
    amount: number,
    unit: 'day' | 'month' | 'year' | 'hour' | 'minute' | 'second'
): Date {
    const d = new Date(date);

    switch (unit) {
        case 'day': d.setDate(d.getDate() + amount); break;
        case 'month': d.setMonth(d.getMonth() + amount); break;
        case 'year': d.setFullYear(d.getFullYear() + amount); break;
        case 'hour': d.setHours(d.getHours() + amount); break;
        case 'minute': d.setMinutes(d.getMinutes() + amount); break;
        case 'second': d.setSeconds(d.getSeconds() + amount); break;
    }

    return d;
}
