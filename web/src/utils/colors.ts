/**
 * Resolves a CSS variable from the DOM and returns its computed string value.
 * Useful for Chart.js which doesn't natively support CSS variables.
 * @param variableName The CSS variable name, e.g., '--color-primary'
 * @returns The resolved color string (hex, rgb, etc.)
 */
export const getCSSVariableValue = (variableName: string): string => {
  if (typeof window === 'undefined') return '';
  return getComputedStyle(document.documentElement)
    .getPropertyValue(variableName)
    .trim();
};
