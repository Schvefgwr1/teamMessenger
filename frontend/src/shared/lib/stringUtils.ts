/**
 * Делает первую букву строки заглавной
 * @param str - строка для преобразования
 * @returns строка с заглавной первой буквой
 * @example
 * capitalizeFirst('hello world') // 'Hello world'
 * capitalizeFirst('HELLO') // 'HELLO'
 */
export function capitalizeFirst(str: string): string {
  if (!str) return str;
  return str.charAt(0).toUpperCase() + str.slice(1);
}

