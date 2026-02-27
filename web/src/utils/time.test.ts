import { describe, test, expect, vi, beforeEach, afterEach } from 'vitest';
import {
  convertTimestampToDate,
  formatTimestamp,
  formatTimestampDate,
  formatRelativeTime,
  isTimestampExpired,
  getRemainingTime,
} from './time';
import { Timestamp } from '@bufbuild/protobuf';

describe('Time Utilities', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe('convertTimestampToDate', () => {
    test('should convert Timestamp to Date', () => {
      const testDate = new Date('2024-01-15T10:30:00Z');
      const timestamp = Timestamp.fromDate(testDate);
      const result = convertTimestampToDate(timestamp);

      expect(result).toBeInstanceOf(Date);
      expect(result?.getTime()).toBe(testDate.getTime());
    });

    test('should return null for undefined timestamp', () => {
      const result = convertTimestampToDate(undefined);
      expect(result).toBeNull();
    });
  });

  describe('formatTimestamp', () => {
    test('should format timestamp correctly', () => {
      const timestamp = Timestamp.fromDate(new Date('2024-01-15T10:30:45Z'));
      const result = formatTimestamp(timestamp);

      expect(result).toMatch(/2024-01-15/);
      expect(result).toMatch(/\d{2}:\d{2}:\d{2}/);
    });

    test('should return "-" for undefined timestamp', () => {
      const result = formatTimestamp(undefined);
      expect(result).toBe('-');
    });

    test('should format with correct padding', () => {
      const timestamp = Timestamp.fromDate(new Date('2024-01-05T08:05:09Z'));
      const result = formatTimestamp(timestamp);

      expect(result).toMatch(/2024-01-05/);
      expect(result).toMatch(/08:05:09/);
    });
  });

  describe('formatTimestampDate', () => {
    test('should format date only', () => {
      const timestamp = Timestamp.fromDate(new Date('2024-01-15T10:30:45Z'));
      const result = formatTimestampDate(timestamp);

      expect(result).toBe('2024-01-15');
    });

    test('should return "-" for undefined timestamp', () => {
      const result = formatTimestampDate(undefined);
      expect(result).toBe('-');
    });
  });

  describe('formatRelativeTime', () => {
    test('should return "刚刚" for recent timestamp', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-15T10:29:30Z'));
      const result = formatRelativeTime(timestamp);

      expect(result).toBe('刚刚');
    });

    test('should return minutes ago', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-15T10:25:00Z'));
      const result = formatRelativeTime(timestamp);

      expect(result).toBe('5分钟前');
    });

    test('should return hours ago', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-15T07:30:00Z'));
      const result = formatRelativeTime(timestamp);

      expect(result).toBe('3小时前');
    });

    test('should return days ago', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-12T10:30:00Z'));
      const result = formatRelativeTime(timestamp);

      expect(result).toBe('3天前');
    });

    test('should return date for old timestamps', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2023-12-01T10:30:00Z'));
      const result = formatRelativeTime(timestamp);

      expect(result).toBe('2023-12-01');
    });

    test('should return "-" for undefined timestamp', () => {
      const result = formatRelativeTime(undefined);
      expect(result).toBe('-');
    });
  });

  describe('isTimestampExpired', () => {
    test('should return true for expired timestamp', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-14T10:30:00Z'));
      const result = isTimestampExpired(timestamp);

      expect(result).toBe(true);
    });

    test('should return false for future timestamp', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-16T10:30:00Z'));
      const result = isTimestampExpired(timestamp);

      expect(result).toBe(false);
    });

    test('should return true for undefined timestamp', () => {
      const result = isTimestampExpired(undefined);
      expect(result).toBe(true);
    });
  });

  describe('getRemainingTime', () => {
    test('should return remaining time in milliseconds', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-15T11:30:00Z'));
      const result = getRemainingTime(timestamp);

      expect(result).toBe(3600000); // 1 hour in milliseconds
    });

    test('should return 0 for expired timestamp', () => {
      const now = new Date('2024-01-15T10:30:00Z');
      vi.setSystemTime(now);

      const timestamp = Timestamp.fromDate(new Date('2024-01-14T10:30:00Z'));
      const result = getRemainingTime(timestamp);

      expect(result).toBe(0);
    });

    test('should return 0 for undefined timestamp', () => {
      const result = getRemainingTime(undefined);
      expect(result).toBe(0);
    });
  });
});
