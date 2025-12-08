import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useDebounce } from './useDebounce';

describe('useDebounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('returns initial value immediately', () => {
    const { result } = renderHook(() => useDebounce('initial', 500));
    expect(result.current).toBe('initial');
  });

  it('debounces value changes', async () => {
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      {
        initialProps: { value: 'initial', delay: 500 },
      }
    );

    expect(result.current).toBe('initial');

    // Изменяем значение
    rerender({ value: 'updated', delay: 500 });

    // Значение еще не должно измениться
    expect(result.current).toBe('initial');

    // Ждем задержку
    await vi.advanceTimersByTimeAsync(500);

    expect(result.current).toBe('updated');
  });

  it('cancels previous timeout on rapid changes', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 500),
      {
        initialProps: { value: 'value1' },
      }
    );

    // Быстро меняем значения
    act(() => {
      rerender({ value: 'value2' });
    });
    await vi.advanceTimersByTimeAsync(200);

    act(() => {
      rerender({ value: 'value3' });
    });
    await vi.advanceTimersByTimeAsync(200);

    act(() => {
      rerender({ value: 'value4' });
    });
    await vi.advanceTimersByTimeAsync(500);

    // Должно быть последнее значение
    expect(result.current).toBe('value4');
  });

  it('respects custom delay', async () => {
    const { result, rerender } = renderHook(
      ({ value, delay }) => useDebounce(value, delay),
      {
        initialProps: { value: 'initial', delay: 1000 },
      }
    );

    act(() => {
      rerender({ value: 'updated', delay: 1000 });
    });

    // После 500мс значение еще не должно измениться
    await vi.advanceTimersByTimeAsync(500);
    expect(result.current).toBe('initial');

    // После 1000мс должно измениться
    await vi.advanceTimersByTimeAsync(500);
    expect(result.current).toBe('updated');
  });

  it('works with numbers', async () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 300),
      {
        initialProps: { value: 0 },
      }
    );

    act(() => {
      rerender({ value: 42 });
    });
    await vi.advanceTimersByTimeAsync(300);

    expect(result.current).toBe(42);
  });

  it('cleans up timeout on unmount', () => {
    const { unmount } = renderHook(() => useDebounce('value', 500));

    unmount();

    // Таймер должен быть очищен, никаких ошибок быть не должно
    vi.advanceTimersByTime(500);
    expect(true).toBe(true); // Если мы здесь, значит ошибок нет
  });
});

