export interface PackSize {
    size: number;
    id: string;
}

export interface CalculationResult {
    packBreakdown: Array<{
        packSize: number;
        quantity: number;
    }>;
    totalItems: number;
    surplusItems: number;
    totalPacks: number;
    calculationTime: number;
}

export interface OptimizationRule {
    minimizeSurplus: boolean;
    minimizePacks: boolean;
}

export interface ApiResponse<T> {
    data: T;
    success: boolean;
    error?: string;
} 