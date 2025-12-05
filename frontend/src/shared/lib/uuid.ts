/**
 * UUID nil значение (пустой UUID)
 * Соответствует uuid.Nil в Go: 00000000-0000-0000-0000-000000000000
 */
export const UUID_NIL = '00000000-0000-0000-0000-000000000000';

/**
 * Проверяет, является ли UUID пустым (nil)
 * @param uuid - UUID строка для проверки
 * @returns true если UUID пустой или равен UUID_NIL
 */
export function isUuidNil(uuid: string | null | undefined): boolean {
  if (!uuid) return true;
  return uuid.trim().toLowerCase() === UUID_NIL.toLowerCase();
}

/**
 * Проверяет, является ли UUID валидным (не nil и не пустой)
 * @param uuid - UUID строка для проверки
 * @returns true если UUID валидный и не nil
 */
export function isValidUuid(uuid: string | null | undefined): boolean {
  if (!uuid) return false;
  const trimmed = uuid.trim();
  if (trimmed === '') return false;
  return trimmed.toLowerCase() !== UUID_NIL.toLowerCase();
}

