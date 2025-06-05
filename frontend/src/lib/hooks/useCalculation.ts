import { useState, useCallback } from 'react';
import { apiService } from '../api';
import type { CalculationResult } from '../types';

export interface UseCalculationReturn {
    result: CalculationResult | null;
    isLoading: boolean;
    error: string | null;
    calculate: (orderQuantity: number, packSizes: number[]) => Promise<void>;
    clearResult: () => void;
}

export function useCalculation(): UseCalculationReturn {
    const [result, setResult] = useState<CalculationResult | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const calculate = useCallback(async (orderQuantity: number, packSizes: number[]) => {
        try {
            setIsLoading(true);
            setError(null);

            // Validate inputs
            if (orderQuantity <= 0) {
                throw new Error('Order quantity must be greater than 0');
            }

            if (packSizes.length === 0) {
                throw new Error('At least one pack size is required');
            }

            // Call the backend API
            const calculationResult = await apiService.calculateOptimalPacks({
                items: orderQuantity,
                pack_sizes: packSizes,
            });

            setResult(calculationResult);
        } catch (err) {
            const message = err instanceof Error ? err.message : 'Calculation failed';
            setError(message);
            setResult(null);
            throw err;
        } finally {
            setIsLoading(false);
        }
    }, []);

    const clearResult = useCallback(() => {
        setResult(null);
        setError(null);
    }, []);

    return {
        result,
        isLoading,
        error,
        calculate,
        clearResult,
    };
} 