import { useEffect, useState } from 'react';
import { Settings2, ChevronDown, ChevronUp } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { StatusAlert } from '@/components/atoms/StatusAlert';
import { ConfigurationCard } from '@/components/atoms/ConfigurationCard';
import { CreateConfigurationForm } from '@/components/atoms/CreateConfigurationForm';
import { usePackConfigurationStore } from '@/lib/store';
import type { PackConfiguration } from '@/lib/types';

interface PackConfigurationPresetsProps {
    currentPackSizes: number[];
    onConfigurationSelect: (packSizes: number[]) => void;
}

export function PackConfigurationPresets({
    currentPackSizes,
    onConfigurationSelect,
}: PackConfigurationPresetsProps) {
    const {
        configurations,
        selectedConfiguration,
        isLoading,
        error,
        loadConfigurations,
        selectConfiguration,
        createConfiguration,
        deleteConfiguration,
        setAsDefault,
        setError,
    } = usePackConfigurationStore();

    const [isOpen, setIsOpen] = useState(false);

    useEffect(() => {
        // Load configurations when component mounts
        loadConfigurations();
    }, [loadConfigurations]);

    const handleSelectConfiguration = (config: PackConfiguration) => {
        selectConfiguration(config);
        onConfigurationSelect(config.pack_sizes);
        setError(null);
    };

    const handleDeleteConfiguration = async (id: number) => {
        await deleteConfiguration(id);
    };

    const handleSetAsDefault = async (id: number) => {
        await setAsDefault(id);
    };

    const handleCreateConfiguration = async (request: { name: string; pack_sizes: number[] }) => {
        await createConfiguration(request);

    };

    const hasConfigurations = configurations.length > 0;

    return (
        <div className="space-y-3">
            <Button
                className="w-full justify-between h-9"
                onClick={() => setIsOpen(!isOpen)}
            >
                <div className="flex items-center gap-2">
                    <Settings2 className="h-4 w-4 text-slate-500" />
                    <span className="text-sm text-slate-600">
                        {selectedConfiguration ? `Using: ${selectedConfiguration.name}` : 'Pack Size Presets'}
                    </span>
                </div>
                {isOpen ? (
                    <ChevronUp className="h-4 w-4 text-slate-400" />
                ) : (
                    <ChevronDown className="h-4 w-4 text-slate-400" />
                )}
            </Button>

            {isOpen && (
                <div className="space-y-3">
                    {error && (
                        <StatusAlert variant="destructive">
                            {error}
                        </StatusAlert>
                    )}

                    {/* Create new configuration */}
                    <CreateConfigurationForm
                        currentPackSizes={currentPackSizes}
                        onSubmit={handleCreateConfiguration}
                        isLoading={isLoading}
                    />

                    {/* Existing configurations */}
                    {hasConfigurations && (
                        <div className="space-y-2">
                            <h4 className="text-sm font-medium text-slate-700">
                                Saved Presets
                            </h4>
                            <div className="grid gap-2">
                                {configurations.map((config) => (
                                    <ConfigurationCard
                                        key={config.id}
                                        configuration={config}
                                        onSelect={handleSelectConfiguration}
                                        onDelete={handleDeleteConfiguration}
                                        onSetDefault={handleSetAsDefault}
                                        isSelected={selectedConfiguration?.id === config.id}
                                        isLoading={isLoading}
                                    />
                                ))}
                            </div>
                        </div>
                    )}

                    {!hasConfigurations && !isLoading && !error && (
                        <div className="text-center py-4 text-sm text-slate-500">
                            No presets saved yet. Create your first one above!
                        </div>
                    )}
                </div>
            )}
        </div>
    );
} 