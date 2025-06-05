import { forwardRef } from 'react';
import type { LucideIcon } from 'lucide-react';
import { Loader2 } from 'lucide-react';
import { Button, buttonVariants } from '@/components/ui/button';
import type { VariantProps } from 'class-variance-authority';

interface LoadingButtonProps extends
    React.ComponentProps<"button">,
    VariantProps<typeof buttonVariants> {
    isLoading?: boolean;
    icon?: LucideIcon;
    asChild?: boolean;
}

export const LoadingButton = forwardRef<HTMLButtonElement, LoadingButtonProps>(
    ({ isLoading = false, icon: Icon, children, disabled, ...props }, ref) => {
        return (
            <Button
                ref={ref}
                disabled={disabled || isLoading}
                {...props}
            >
                <div className="flex items-center gap-2">
                    {isLoading ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                    ) : Icon ? (
                        <Icon className="h-4 w-4" />
                    ) : null}
                    {children}
                </div>
            </Button>
        );
    }
);

LoadingButton.displayName = 'LoadingButton'; 