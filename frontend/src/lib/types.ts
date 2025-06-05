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

export interface PackConfiguration {
    id: number;
    name: string;
    pack_sizes: number[];
    is_default: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreatePackConfigurationRequest {
    name: string;
    pack_sizes: number[];
}

export interface UpdatePackConfigurationRequest {
    name: string;
    pack_sizes: number[];
    is_default?: boolean;
} 