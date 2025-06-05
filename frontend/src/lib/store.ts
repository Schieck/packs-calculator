import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { apiService } from './api';
import type { PackConfiguration, CreatePackConfigurationRequest } from './types';

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

// Pack Configuration Store
interface PackConfigurationStore {
    // State
    configurations: PackConfiguration[];
    selectedConfiguration: PackConfiguration | null;
    isLoading: boolean;
    error: string | null;

    // Actions
    loadConfigurations: () => Promise<void>;
    loadDefaultConfiguration: () => Promise<void>;
    selectConfiguration: (config: PackConfiguration) => void;
    createConfiguration: (request: CreatePackConfigurationRequest) => Promise<void>;
    deleteConfiguration: (id: number) => Promise<void>;
    setAsDefault: (id: number) => Promise<void>;
    setError: (error: string | null) => void;
    setLoading: (loading: boolean) => void;
}

export const usePackConfigurationStore = create<PackConfigurationStore>()((set, get) => ({
    // Initial state
    configurations: [],
    selectedConfiguration: null,
    isLoading: false,
    error: null,

    // Actions
    loadConfigurations: async () => {
        try {
            set({ isLoading: true, error: null });
            const configurations = await apiService.getAllPackConfigurations();
            set({ configurations, isLoading: false });
        } catch (error) {
            console.error('Failed to load configurations:', error);
            set({
                error: 'Failed to load pack configurations',
                isLoading: false
            });
        }
    },

    loadDefaultConfiguration: async () => {
        try {
            set({ isLoading: true, error: null });
            const defaultConfig = await apiService.getDefaultPackConfiguration();
            if (defaultConfig) {
                set({
                    selectedConfiguration: defaultConfig,
                    isLoading: false
                });
            } else {
                set({ isLoading: false });
            }
        } catch (error) {
            console.error('Failed to load default configuration:', error);
            set({
                error: 'Failed to load default configuration',
                isLoading: false
            });
        }
    },

    selectConfiguration: (config: PackConfiguration) => {
        set({ selectedConfiguration: config, error: null });
    },

    createConfiguration: async (request: CreatePackConfigurationRequest) => {
        try {
            set({ isLoading: true, error: null });
            const newConfig = await apiService.createPackConfiguration(request);
            const { configurations } = get();
            set({
                configurations: [...configurations, newConfig],
                isLoading: false
            });
        } catch (error) {
            console.error('Failed to create configuration:', error);
            set({
                error: 'Failed to create pack configuration',
                isLoading: false
            });
            throw error;
        }
    },

    deleteConfiguration: async (id: number) => {
        try {
            set({ isLoading: true, error: null });
            await apiService.deletePackConfiguration(id);
            const { configurations, selectedConfiguration } = get();
            const updatedConfigurations = configurations.filter(c => c.id !== id);
            const updatedSelected = selectedConfiguration?.id === id ? null : selectedConfiguration;
            set({
                configurations: updatedConfigurations,
                selectedConfiguration: updatedSelected,
                isLoading: false
            });
        } catch (error) {
            console.error('Failed to delete configuration:', error);
            set({
                error: 'Failed to delete pack configuration',
                isLoading: false
            });
            throw error;
        }
    },

    setAsDefault: async (id: number) => {
        try {
            set({ isLoading: true, error: null });
            const updatedConfig = await apiService.setDefaultPackConfiguration(id);
            const { configurations } = get();

            const updatedConfigurations = configurations.map(config => ({
                ...config,
                is_default: config.id === id
            }));

            set({
                configurations: updatedConfigurations,
                selectedConfiguration: updatedConfig,
                isLoading: false
            });
        } catch (error) {
            console.error('Failed to set default configuration:', error);
            set({
                error: 'Failed to set default configuration',
                isLoading: false
            });
            throw error;
        }
    },

    setError: (error: string | null) => {
        set({ error });
    },

    setLoading: (isLoading: boolean) => {
        set({ isLoading });
    },
})); 