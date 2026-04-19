/**
 * Object utility functions
 * Common operations for object manipulation
 */

type UnknownRecord = Record<PropertyKey, unknown>;

/**
 * Deep clones an object
 * Uses the native structuredClone where available, with a basic fallback
 * @param obj - Object to clone
 */
export function deepClone<T>(obj: T): T {
  if (typeof structuredClone === 'function') {
    return structuredClone(obj);
  }

  // Basic fallback for environments without structuredClone
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }
  
  if (Array.isArray(obj)) {
    return obj.map(item => deepClone(item)) as unknown as T;
  }
  
  const source = obj as UnknownRecord;
  const clonedObj: UnknownRecord = {};

  Object.keys(source).forEach((key) => {
    clonedObj[key] = deepClone(source[key]);
  });
  
  return clonedObj as T;
}

/**
 * Safely gets a nested property from an object using a path string
 * @param obj - The object to get the value from
 * @param path - Path to the property (e.g., 'user.address.city')
 * @param defaultValue - Default value if path doesn't exist
 */
export function getNestedValue<T, D = undefined>(
  obj: Record<string, unknown> | null | undefined,
  path: string,
  defaultValue?: D
): T | D {
  if (!obj) return defaultValue as D;
  
  const keys = path.split('.');
  let result: unknown = obj;
  
  for (const key of keys) {
    if (
      result === undefined ||
      result === null ||
      (typeof result !== 'object' && typeof result !== 'function')
    ) {
      return defaultValue as D;
    }

    result = Reflect.get(result as object, key);
  }
  
  return (result === undefined) ? (defaultValue as D) : (result as T);
}

/**
 * Creates a new object with only the specified keys
 * @param obj - Original object
 * @param keys - Keys to include in the new object
 */
export function pick<T extends object, K extends keyof T>(
  obj: T, 
  keys: K[]
): Pick<T, K> {
  const result = {} as Pick<T, K>;
  keys.forEach(key => {
    if (obj && Object.prototype.hasOwnProperty.call(obj, key)) {
      result[key] = obj[key];
    }
  });
  return result;
}

/**
 * Creates a new object excluding the specified keys
 * @param obj - Original object
 * @param keys - Keys to exclude from the new object
 */
export function omit<T extends object, K extends keyof T>(
  obj: T, 
  keys: K[]
): Omit<T, K> {
  const result = { ...obj };
  const mutableResult = result as UnknownRecord;

  keys.forEach(key => {
    delete mutableResult[key];
  });

  return result;
}

/**
 * Checks if an object is empty
 * @param obj - Object to check
 */
export function isEmptyObject(obj: object | null | undefined): boolean {
  if (!obj) return true;
  return Object.keys(obj).length === 0;
}

/**
 * Merges objects deeply
 * @param target - Target object
 * @param sources - Source objects
 */
export function deepMerge<T extends object>(target: T, ...sources: object[]): T {
  if (!sources.length) return target;
  const source = sources.shift();

  if (isObject(target) && isObject(source)) {
    const targetRecord = target as UnknownRecord;

    Object.keys(source).forEach(key => {
      const sourceValue = source[key];
      const targetValue = targetRecord[key];

      if (isObject(sourceValue)) {
        if (!isObject(targetValue)) {
          targetRecord[key] = deepClone(sourceValue);
        } else {
          targetRecord[key] = deepMerge(targetValue, sourceValue);
        }
      } else {
        targetRecord[key] = sourceValue;
      }
    });
  }

  return deepMerge(target, ...sources);
}

/**
 * Deep freezes an object
 * @param obj - Object to freeze
 */
export function deepFreeze<T extends object>(obj: T): T {
  Object.freeze(obj);
  Object.getOwnPropertyNames(obj).forEach(prop => {
    const value = Reflect.get(obj as object, prop);
    if (
      value !== null &&
      (typeof value === 'object' || typeof value === 'function') &&
      !Object.isFrozen(value)
    ) {
      deepFreeze(value);
    }
  });
  return obj;
}

/**
 * Helper: Checks if value is a plain object
 */
function isObject(item: unknown): item is Record<string, unknown> {
  return item !== null && typeof item === 'object' && !Array.isArray(item);
}
