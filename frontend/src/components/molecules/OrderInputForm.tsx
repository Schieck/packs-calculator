import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { Calculator } from 'lucide-react';

import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LoadingButton } from '@/components/atoms/LoadingButton';
import { orderCalculationSchema, type OrderCalculationFormData } from '@/lib/schemas';

interface OrderInputFormProps {
    onSubmit: (data: OrderCalculationFormData) => Promise<void>;
    isLoading: boolean;
    isDisabled: boolean;
    onKeyDown?: (e: React.KeyboardEvent) => void;
}

export function OrderInputForm({
    onSubmit,
    isLoading,
    isDisabled,
    onKeyDown
}: OrderInputFormProps) {
    const {
        register,
        handleSubmit,
        formState: { errors, isSubmitting },
    } = useForm<OrderCalculationFormData>({
        resolver: zodResolver(orderCalculationSchema),
    });

    return (
        <div className="space-y-2">
            <Label htmlFor="orderQuantity" className="text-sm font-medium">
                Order Quantity
            </Label>
            <form onSubmit={handleSubmit(onSubmit)} className="flex gap-2">
                <div className="flex-1">
                    <Input
                        {...register('orderQuantity')}
                        id="orderQuantity"
                        type="number"
                        placeholder="Enter order quantity (e.g., 1,234)"
                        className="w-full text-lg h-12"
                        min="1"
                        max="1000000000"
                        autoComplete="off"
                        onKeyDown={onKeyDown}
                        disabled={isDisabled}
                    />
                    {errors.orderQuantity && (
                        <p className="text-sm text-destructive mt-1">
                            {errors.orderQuantity.message}
                        </p>
                    )}
                </div>
                <LoadingButton
                    type="submit"
                    disabled={isDisabled}
                    className="h-12 px-6"
                    size="lg"
                    isLoading={isLoading || isSubmitting}
                    icon={Calculator}
                >
                    Calculate
                </LoadingButton>
            </form>
            <p className="text-xs text-muted-foreground">
                Press Enter to calculate, Escape to clear results
            </p>
        </div>
    );
} 