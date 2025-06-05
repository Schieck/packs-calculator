import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Plus } from 'lucide-react';

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LoadingButton } from '@/components/atoms/LoadingButton';
import { addPackSizeSchema, type AddPackSizeFormData } from '@/lib/schemas';

interface AddPackSizeFormProps {
    onSubmit: (data: AddPackSizeFormData) => Promise<void>;
    isLoading: boolean;
}

export function AddPackSizeForm({ onSubmit, isLoading }: AddPackSizeFormProps) {
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors, isSubmitting },
    } = useForm<AddPackSizeFormData>({
        resolver: zodResolver(addPackSizeSchema),
    });

    const handleFormSubmit = async (data: AddPackSizeFormData) => {
        await onSubmit(data);
        reset();
    };

    return (
        <div className="space-y-3">
            <Label className="text-sm font-medium">Add New Pack Size</Label>
            <form onSubmit={handleSubmit(handleFormSubmit)} className="flex gap-2">
                <div className="flex-1">
                    <Input
                        {...register('newPackSize')}
                        type="number"
                        placeholder="Enter pack size (e.g., 750)"
                        className="w-full"
                        min="1"
                        max="10000000"
                        autoComplete="off"
                    />
                    {errors.newPackSize && (
                        <p className="text-sm text-destructive mt-1">
                            {errors.newPackSize.message}
                        </p>
                    )}
                </div>
                <LoadingButton
                    type="submit"
                    disabled={isSubmitting}
                    isLoading={isLoading}
                    icon={Plus}
                >
                    Add
                </LoadingButton>
            </form>
        </div>
    );
} 