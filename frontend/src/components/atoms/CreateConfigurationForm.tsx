import { useState } from 'react';
import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { CreatePackConfigurationRequest } from '@/lib/types';

interface CreateConfigurationFormProps {
    currentPackSizes: number[];
    onSubmit: (request: CreatePackConfigurationRequest) => Promise<void>;
    isLoading?: boolean;
}

export function CreateConfigurationForm({
    currentPackSizes,
    onSubmit,
    isLoading = false,
}: CreateConfigurationFormProps) {
    const [name, setName] = useState('');
    const [isExpanded, setIsExpanded] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!name.trim()) {
            setError('Configuration name is required');
            return;
        }

        if (currentPackSizes.length === 0) {
            setError('At least one pack size is required');
            return;
        }

        try {
            setError(null);
            await onSubmit({
                name: name.trim(),
                pack_sizes: currentPackSizes,
            });
            setName('');
            setIsExpanded(false);
        } catch {
            setError('Failed to create configuration');
        }
    };

    if (!isExpanded) {
        return (
            <Button
                size="sm"
                onClick={() => setIsExpanded(true)}
                className="w-full"
            >
                <Plus className="h-4 w-4 mr-2" />
                Save Current as Preset
            </Button>
        );
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-3 p-3 border rounded-lg bg-slate-50">
            <div className="space-y-2">
                <Label htmlFor="config-name" className="text-sm">
                    Preset Name
                </Label>
                <Input
                    id="config-name"
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="e.g., 'Small Items', 'Bulk Orders'"
                    className="h-8"
                    disabled={isLoading}
                />
            </div>

            {error && (
                <p className="text-xs text-red-600">{error}</p>
            )}

            <div className="flex gap-2">
                <Button
                    type="submit"
                    size="sm"
                    disabled={isLoading || !name.trim()}
                    className="flex-1"
                >
                    {isLoading ? 'Saving...' : 'Save Preset'}
                </Button>
                <Button
                    type="button"
                    size="sm"
                    onClick={() => {
                        setIsExpanded(false);
                        setName('');
                        setError(null);
                    }}
                    disabled={isLoading}
                >
                    Cancel
                </Button>
            </div>
        </form>
    );
} 