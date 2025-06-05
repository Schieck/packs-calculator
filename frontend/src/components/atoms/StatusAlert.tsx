import type { LucideIcon } from 'lucide-react';
import { Alert, AlertDescription } from '@/components/ui/alert';

interface StatusAlertProps {
    children: React.ReactNode;
    variant?: 'default' | 'destructive';
    icon?: LucideIcon;
    className?: string;
    onClick?: () => void;
}

export function StatusAlert({
    children,
    variant = 'default',
    icon: Icon,
    className,
    onClick
}: StatusAlertProps) {
    return (
        <Alert
            variant={variant}
            className={`${className || ''} ${onClick ? 'cursor-pointer hover:bg-opacity-80 transition-colors' : ''}`}
            onClick={onClick}
        >
            {Icon && <Icon className="h-4 w-4" />}
            <AlertDescription>{children}</AlertDescription>
        </Alert>
    );
} 