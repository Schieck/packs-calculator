import type { CalculationResult } from './types';

export interface CalculateRequest {
    items: number;
    pack_sizes: number[];
}

export interface CalculateResponse {
    allocation: Record<string, number>;
    total_packs: number;
    total_items: number;
    surplus: number;
}

export interface AuthRequest {
    secret: string;
}

export interface AuthResponse {
    token: string;
}

class ApiService {
    private baseUrl: string;
    private cachedToken: string | null = null;
    private tokenExpiryTime: number | null = null;

    constructor() {
        // In production, this would be from environment variables
        this.baseUrl = 'http://localhost:8080/api/v1';
    }

    private async getAuthToken(): Promise<string> {
        // Check if we have a valid cached token (assuming 1 hour expiry minus 5 minutes buffer)
        if (this.cachedToken && this.tokenExpiryTime && Date.now() < this.tokenExpiryTime) {
            return this.cachedToken;
        }

        try {
            const authRequest: AuthRequest = {
                secret: import.meta.env.VITE_AUTH_SECRET || 'your-auth-secret-change-in-production'
            };

            const response = await fetch(`${this.baseUrl}/auth/token`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(authRequest),
            });

            if (!response.ok) {
                throw new Error(`Authentication failed: ${response.status} ${response.statusText}`);
            }

            const data: AuthResponse = await response.json();

            // Cache the token with expiry (assuming 1 hour token lifetime minus 5 minutes buffer)
            this.cachedToken = data.token;
            this.tokenExpiryTime = Date.now() + (55 * 60 * 1000); // 55 minutes

            return data.token;
        } catch (error) {
            console.error('Authentication error:', error);
            throw new Error('Failed to authenticate with the API');
        }
    }

    async calculateOptimalPacks(request: CalculateRequest): Promise<CalculationResult> {
        const startTime = performance.now();

        try {
            // Get authentication token
            const token = await this.getAuthToken();

            const response = await fetch(`${this.baseUrl}/calculate`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify(request),
            });

            if (!response.ok) {
                // If authentication fails, clear cached token and retry once
                if (response.status === 401 && this.cachedToken) {
                    this.cachedToken = null;
                    this.tokenExpiryTime = null;

                    // Retry with fresh token
                    const newToken = await this.getAuthToken();
                    const retryResponse = await fetch(`${this.baseUrl}/calculate`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `Bearer ${newToken}`
                        },
                        body: JSON.stringify(request),
                    });

                    if (!retryResponse.ok) {
                        throw new Error(`API request failed: ${retryResponse.status} ${retryResponse.statusText}`);
                    }

                    const retryData: CalculateResponse = await retryResponse.json();
                    return this.transformResponse(retryData, performance.now() - startTime);
                }

                throw new Error(`API request failed: ${response.status} ${response.statusText}`);
            }

            const data: CalculateResponse = await response.json();
            return this.transformResponse(data, performance.now() - startTime);
        } catch (error) {
            console.error('API calculation error:', error);
            throw error;
        }
    }

    private transformResponse(data: CalculateResponse, calculationTime: number): CalculationResult {
        // Transform API response to our frontend types
        const packBreakdown = Object.entries(data.allocation).map(([packSize, quantity]) => ({
            packSize: parseInt(packSize),
            quantity: quantity,
        })).sort((a, b) => b.packSize - a.packSize); // Sort by pack size descending

        return {
            packBreakdown,
            totalItems: data.total_items,
            surplusItems: data.surplus,
            totalPacks: data.total_packs,
            calculationTime,
        };
    }

    // Method to manually clear cached token (useful for logout or token refresh scenarios)
    clearAuthToken(): void {
        this.cachedToken = null;
        this.tokenExpiryTime = null;
    }

    // Mock localStorage-based pack management for development
    private PACK_SIZES_KEY = 'pack-calculator-sizes';
    private DEFAULT_PACK_SIZES = [250, 500, 1000, 2000, 5000];

    getPackSizes(): number[] {
        try {
            const stored = localStorage.getItem(this.PACK_SIZES_KEY);
            return stored ? JSON.parse(stored) : this.DEFAULT_PACK_SIZES;
        } catch {
            return this.DEFAULT_PACK_SIZES;
        }
    }

    savePackSizes(sizes: number[]): void {
        try {
            localStorage.setItem(this.PACK_SIZES_KEY, JSON.stringify(sizes));
        } catch (error) {
            console.error('Failed to save pack sizes:', error);
        }
    }
}

export const apiService = new ApiService(); 