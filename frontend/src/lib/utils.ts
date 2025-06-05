import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Formatting utilities
export const formatNumber = (num: number): string => num.toLocaleString();

export const formatTime = (ms: number): string => {
  if (ms < 1) return `${ms.toFixed(2)}ms`;
  if (ms < 1000) return `${Math.round(ms)}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
};

// Validation utilities
export const isValidPackSize = (size: number): boolean => size >= 1 && size <= 10000000;

export const isDuplicatePackSize = (size: number, existingSizes: number[]): boolean =>
  existingSizes.includes(size);
