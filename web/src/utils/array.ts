/**
 * Array utility functions
 * Common operations for array manipulation
 */

/**
 * Groups array items by a specified key
 * @param array - Array to group
 * @param key - Object key to group by
 */
export function groupBy<T>(array: T[], key: keyof T): Record<string, T[]> {
  return array.reduce((result, item) => {
    const groupKey = String(item[key]);
    if (!result[groupKey]) {
      result[groupKey] = [];
    }
    result[groupKey].push(item);
    return result;
  }, {} as Record<string, T[]>);
}

/**
 * Creates a unique array from an array with possible duplicates
 * @param array - Array with possible duplicates
 */
export function uniqueArray<T>(array: T[]): T[] {
  return [...new Set(array)];
}

/**
 * Sorts an array by a specified object key
 * @param array - Array to sort
 * @param key - Object key to sort by
 * @param direction - Sort direction (asc or desc)
 */
export function sortByKey<T>(
  array: T[], 
  key: keyof T, 
  direction: 'asc' | 'desc' = 'asc'
): T[] {
  return [...array].sort((a, b) => {
    if (a[key] < b[key]) return direction === 'asc' ? -1 : 1;
    if (a[key] > b[key]) return direction === 'asc' ? 1 : -1;
    return 0;
  });
}

/**
 * Chunks an array into smaller arrays of specified size
 * @param array - Array to chunk
 * @param size - Size of each chunk
 */
export function chunkArray<T>(array: T[], size: number): T[][] {
  return array.reduce((result, item, index) => {
    const chunkIndex = Math.floor(index / size);
    
    if (!result[chunkIndex]) {
      result[chunkIndex] = [];
    }
    
    result[chunkIndex].push(item);
    
    return result;
  }, [] as T[][]);
}

/**
 * Gets the intersection of two arrays
 * @param arrayA - First array
 * @param arrayB - Second array
 */
export function arrayIntersection<T>(arrayA: T[], arrayB: T[]): T[] {
  return arrayA.filter(item => arrayB.includes(item));
}

/**
 * Gets the difference between two arrays
 * @param arrayA - Array to compare from
 * @param arrayB - Array to compare against
 */
export function arrayDifference<T>(arrayA: T[], arrayB: T[]): T[] {
  return arrayA.filter(item => !arrayB.includes(item));
}
