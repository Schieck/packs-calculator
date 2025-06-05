import type { LucideIcon } from 'lucide-react';
import { formatNumber } from '@/lib/utils';

interface MetricCardProps {
    title: string;
    value: number;
    subtitle?: string;
    icon: LucideIcon;
    variant?: 'blue' | 'orange' | 'green' | 'gray';
}

const variantStyles = {
    blue: 'bg-blue-50 border-blue-200 text-blue-600',
    orange: 'bg-orange-50 border-orange-200 text-orange-600',
    green: 'bg-green-50 border-green-200 text-green-600',
    gray: 'bg-gray-50 border-gray-200 text-gray-600',
};

const titleStyles = {
    blue: 'text-blue-900',
    orange: 'text-orange-900',
    green: 'text-green-900',
    gray: 'text-gray-900',
};

const valueStyles = {
    blue: 'text-blue-700',
    orange: 'text-orange-700',
    green: 'text-green-700',
    gray: 'text-gray-700',
};

const subtitleStyles = {
    blue: 'text-blue-600',
    orange: 'text-orange-600',
    green: 'text-green-600',
    gray: 'text-gray-600',
};

export function MetricCard({
    title,
    value,
    subtitle,
    icon: Icon,
    variant = 'blue'
}: MetricCardProps) {
    return (
        <div className={`border rounded-lg p-4 ${variantStyles[variant]}`}>
            <div className="flex items-center gap-2 mb-2">
                <Icon className="h-4 w-4" />
                <span className={`text-sm font-medium ${titleStyles[variant]}`}>
                    {title}
                </span>
            </div>
            <div className={`text-2xl font-bold ${valueStyles[variant]}`}>
                {formatNumber(value)}
            </div>
            {subtitle && (
                <div className={`text-xs mt-1 ${subtitleStyles[variant]}`}>
                    {subtitle}
                </div>
            )}
        </div>
    );
} 