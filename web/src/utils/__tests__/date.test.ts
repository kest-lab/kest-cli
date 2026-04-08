import { describe, it, expect } from 'vitest';
import { formatDate, getRelativeTime, isToday, addTime } from '../date';

describe('date utils', () => {
  describe('formatDate', () => {
    const testDate = new Date('2023-10-25T14:30:00');

    it('should format with default YYYY-MM-DD', () => {
      expect(formatDate(testDate)).toBe('2023-10-25');
    });

    it('should format with custom string format', () => {
      expect(formatDate(testDate, 'YYYY/MM/DD HH:mm')).toBe('2023/10/25 14:30');
    });

    it('should use Intl for options', () => {
      const formatted = formatDate(testDate, { month: 'long', year: 'numeric' }, 'en-US');
      expect(formatted).toContain('October');
      expect(formatted).toContain('2023');
    });
  });

  describe('getRelativeTime', () => {
    const now = new Date();

    it('should format seconds ago', () => {
      const past = new Date(now.getTime() - 10000);
      expect(getRelativeTime(past, now)).toContain('seconds ago');
    });

    it('should format hours ago', () => {
      const past = new Date(now.getTime() - 3600 * 2000);
      expect(getRelativeTime(past, now)).toContain('2 hours ago');
    });
  });

  describe('isToday', () => {
    it('should return true for today', () => {
      expect(isToday(new Date())).toBe(true);
    });

    it('should return false for yesterday', () => {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);
      expect(isToday(yesterday)).toBe(false);
    });
  });

  describe('addTime', () => {
    it('should add days correctly', () => {
      const d = new Date('2023-01-01');
      const result = addTime(d, 5, 'day');
      expect(result.getDate()).toBe(6);
    });
  });
});
