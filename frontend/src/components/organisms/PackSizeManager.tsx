import { useState } from 'react';
import { Package } from 'lucide-react';

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';

import { StatusAlert } from '@/components/atoms/StatusAlert';
import { PackSizeList } from '@/components/molecules/PackSizeList';
import { AddPackSizeForm } from '@/components/molecules/AddPackSizeForm';
import { BulkActions } from '@/components/molecules/BulkActions';

import { usePackSizesStore } from '@/lib/store';
import { type AddPackSizeFormData } from '@/lib/schemas';

export function PackSizeManager() {
    const {
        packSizes,
        error,
        addPackSize,
        removePackSize,
        resetToDefaults,
        clearAll,
        setError,
    } = usePackSizesStore();

    const [isAddingPack, setIsAddingPack] = useState(false);

    const handleAddPackSize = async (data: AddPackSizeFormData) => {
        try {
            setIsAddingPack(true);
            setError(null);
            const size = parseInt(data.newPackSize, 10);
            addPackSize(size);
        } catch (err) {
            console.error('Failed to add pack size:', err);
        } finally {
            setIsAddingPack(false);
        }
    };

    const handleRemovePackSize = (size: number) => {
        setError(null);
        removePackSize(size);
    };

    const handleResetToDefaults = () => {
        setError(null);
        resetToDefaults();
    };

    const handleClearAll = () => {
        setError(null);
        clearAll();
    };

    return (
        <Card className="w-full">
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Package className="h-5 w-5 text-blue-600" />
                    Pack Size Management
                </CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
                {/* Error Display */}
                {error && (
                    <StatusAlert variant="destructive">
                        {error}
                    </StatusAlert>
                )}

                {/* Current Pack Sizes */}
                <PackSizeList
                    packSizes={packSizes}
                    onRemove={handleRemovePackSize}
                />

                <Separator />

                {/* Add New Pack Size */}
                <AddPackSizeForm
                    onSubmit={handleAddPackSize}
                    isLoading={isAddingPack}
                />

                <Separator />

                {/* Bulk Actions */}
                <BulkActions
                    onResetToDefaults={handleResetToDefaults}
                    onClearAll={handleClearAll}
                />
            </CardContent>
        </Card>
    );
} 