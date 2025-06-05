import { z } from 'zod';

export const packSizeSchema = z.object({
    size: z
        .number()
        .int('Pack size must be a whole number')
        .min(1, 'Pack size must be at least 1')
        .max(10000000, 'Pack size must be reasonable (max 10,000,000)'),
});

export const addPackSizeSchema = z.object({
    newPackSize: z
        .string()
        .min(1, 'Pack size is required')
        .refine(
            (val) => {
                const num = parseInt(val, 10);
                return !isNaN(num) && num > 0 && num <= 10000000;
            },
            'Must be a valid number between 1 and 10,000,000'
        ),
});

export const orderCalculationSchema = z.object({
    orderQuantity: z
        .string()
        .min(1, 'Order quantity is required')
        .refine(
            (val) => {
                const num = parseInt(val, 10);
                return !isNaN(num) && num > 0 && num <= 1000000000;
            },
            'Must be a valid number between 1 and 1,000,000,000'
        ),
});

export type PackSizeFormData = z.infer<typeof packSizeSchema>;
export type AddPackSizeFormData = z.infer<typeof addPackSizeSchema>;
export type OrderCalculationFormData = z.infer<typeof orderCalculationSchema>; 