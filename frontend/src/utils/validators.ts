/**
 * Validate URL format
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url);
    return true;
  } catch {
    return false;
  }
}

/**
 * Validate that at least one of the fields is provided
 */
export function hasAtLeastOne(
  values: (string | undefined)[],
  minLength = 1
): boolean {
  return values.some((value) => value && value.length >= minLength);
}

/**
 * Validate required field
 */
export function isRequired(value: string | undefined): boolean {
  return !!value && value.trim().length > 0;
}

/**
 * Validate minimum length
 */
export function minLength(value: string, min: number): boolean {
  return value.length >= min;
}

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Validate number within range
 */
export function isInRange(
  value: number,
  min: number,
  max: number
): boolean {
  return value >= min && value <= max;
}
