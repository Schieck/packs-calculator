import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface PackSizesStore {
    // State
    packSizes: number[];
    isLoading: boolean;
    error: string | null;

    // Actions
    addPackSize: (size: number) => void;
    removePackSize: (size: number) => void;
    resetToDefaults: () => void;
    clearAll: () => void;
    setError: (error: string | null) => void;
    setLoading: (loading: boolean) => void;
}

const DEFAULT_PACK_SIZES = [250, 500, 1000, 2000, 5000];

export const usePackSizesStore = create<PackSizesStore>()(
    persist(
        (set, get) => ({
            // Initial state
            packSizes: DEFAULT_PACK_SIZES,
            isLoading: false,
            error: null,

            // Actions
            addPackSize: (size: number) => {
                const { packSizes } = get();

                // Validation
                if (size < 1 || size > 10000000) {
                    set({ error: 'Pack size must be between 1 and 10,000,000' });
                    return;
                }

                if (packSizes.includes(size)) {
                    set({ error: 'Pack size already exists' });
                    return;
                }

                // Add and sort
                const newSizes = [...packSizes, size].sort((a, b) => a - b);
                set({ packSizes: newSizes, error: null });
            },

            removePackSize: (size: number) => {
                const { packSizes } = get();
                const newSizes = packSizes.filter(s => s !== size);
                set({ packSizes: newSizes, error: null });
            },

            resetToDefaults: () => {
                set({ packSizes: DEFAULT_PACK_SIZES, error: null });
            },

            clearAll: () => {
                set({ packSizes: [], error: null });
            },

            setError: (error: string | null) => {
                set({ error });
            },

            setLoading: (isLoading: boolean) => {
                set({ isLoading });
            },
        }),
        {
            name: 'pack-sizes-storage', // Key in localStorage
            partialize: (state) => ({ packSizes: state.packSizes }), // Only persist packSizes
        }
    )
);

// Convenience hook for components that only need pack sizes
export const usePackSizes = () => usePackSizesStore((state) => state.packSizes); 