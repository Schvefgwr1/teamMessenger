import { format, formatDistanceToNow, isToday, isYesterday, parseISO } from 'date-fns';
import { ru } from 'date-fns/locale';

/**
 * Форматировать дату в человекочитаемый формат
 */
export function formatDate(dateString: string, formatStr = 'dd.MM.yyyy'): string {
  try {
    const date = parseISO(dateString);
    return format(date, formatStr, { locale: ru });
  } catch {
    return dateString;
  }
}

/**
 * Форматировать дату и время
 */
export function formatDateTime(dateString: string): string {
  return formatDate(dateString, 'dd.MM.yyyy HH:mm');
}

/**
 * Форматировать время
 */
export function formatTime(dateString: string): string {
  return formatDate(dateString, 'HH:mm');
}

/**
 * Относительное время (например, "5 минут назад")
 */
export function formatRelativeTime(dateString: string): string {
  try {
    const date = parseISO(dateString);
    return formatDistanceToNow(date, { addSuffix: true, locale: ru });
  } catch {
    return dateString;
  }
}

/**
 * Умное форматирование даты для разделителя в чате
 * - Сегодня: "Сегодня"
 * - Вчера: "Вчера"
 * - Этот год: день и месяц (например, "15 дек")
 * - Другой год: полная дата (например, "15 дек 2023")
 */
export function formatChatDate(dateString: string): string {
  try {
    const date = parseISO(dateString);
    
    if (isToday(date)) {
      return 'Сегодня';
    }
    
    if (isYesterday(date)) {
      return 'Вчера';
    }
    
    const currentYear = new Date().getFullYear();
    const dateYear = date.getFullYear();
    
    if (currentYear === dateYear) {
      return format(date, 'd MMM', { locale: ru });
    }
    
    return format(date, 'd MMM yyyy', { locale: ru });
  } catch {
    return dateString;
  }
}

/**
 * Умное форматирование даты для сообщения
 */
export function formatMessageTime(dateString: string): string {
  try {
    const date = parseISO(dateString);
    return format(date, 'HH:mm', { locale: ru });
  } catch {
    return dateString;
  }
}

