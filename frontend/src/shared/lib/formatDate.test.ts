import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { formatDate, formatDateTime, formatTime, formatRelativeTime, formatChatDate, formatMessageTime } from './formatDate';

describe('formatDate', () => {
  it('formats ISO date string', () => {
    const result = formatDate('2024-01-15T10:30:00Z');
    expect(result).toMatch(/15\.01\.2024/);
  });

  it('formats with custom format', () => {
    const result = formatDate('2024-01-15T10:30:00Z', 'yyyy-MM-dd');
    expect(result).toBe('2024-01-15');
  });

  it('returns original string on invalid date', () => {
    const result = formatDate('invalid-date');
    expect(result).toBe('invalid-date');
  });
});

describe('formatDateTime', () => {
  it('formats date and time', () => {
    const result = formatDateTime('2024-01-15T10:30:00Z');
    expect(result).toMatch(/15\.01\.2024/);
    // Время может отличаться из-за часового пояса, проверяем только наличие времени
    expect(result).toMatch(/\d{2}:\d{2}/);
  });
});

describe('formatTime', () => {
  it('formats time only', () => {
    const result = formatTime('2024-01-15T10:30:00Z');
    // Время может отличаться из-за часового пояса, проверяем только формат
    expect(result).toMatch(/\d{2}:\d{2}/);
  });
});

describe('formatRelativeTime', () => {
  it('formats relative time', () => {
    const now = new Date();
    const past = new Date(now.getTime() - 5 * 60 * 1000); // 5 minutes ago
    const result = formatRelativeTime(past.toISOString());
    expect(result).toContain('назад');
  });
});

describe('formatChatDate', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    // Устанавливаем фиксированную дату для тестов
    vi.setSystemTime(new Date('2024-01-15T12:00:00Z'));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns "Сегодня" for today', () => {
    const today = new Date('2024-01-15T10:00:00Z');
    const result = formatChatDate(today.toISOString());
    expect(result).toBe('Сегодня');
  });

  it('returns "Вчера" for yesterday', () => {
    const yesterday = new Date('2024-01-14T10:00:00Z');
    const result = formatChatDate(yesterday.toISOString());
    expect(result).toBe('Вчера');
  });

  it('formats date for same year', () => {
    const date = new Date('2024-03-20T10:00:00Z');
    const result = formatChatDate(date.toISOString());
    expect(result).toMatch(/20/);
    expect(result).toMatch(/мар/);
  });

  it('formats full date for different year', () => {
    const date = new Date('2023-03-20T10:00:00Z');
    const result = formatChatDate(date.toISOString());
    expect(result).toMatch(/2023/);
  });
});

describe('formatMessageTime', () => {
  it('formats time for message', () => {
    const result = formatMessageTime('2024-01-15T10:30:00Z');
    // Время может отличаться из-за часового пояса, проверяем только формат
    expect(result).toMatch(/\d{2}:\d{2}/);
  });

  it('returns original string on invalid date', () => {
    const result = formatMessageTime('invalid');
    expect(result).toBe('invalid');
  });
});

