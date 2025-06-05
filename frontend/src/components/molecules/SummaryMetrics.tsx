import { Package2, TrendingDown } from 'lucide-react';
import { MetricCard } from '@/components/atoms/MetricCard';
import type { CalculationResult } from '@/lib/types';

interface SummaryMetricsProps {
    result: CalculationResult;
}

export function SummaryMetrics({ result }: SummaryMetricsProps) {
    return (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <MetricCard
                title="Total Items Shipped"
                value={result.totalItems}
                icon={Package2}
                variant="blue"
            />

            <MetricCard
                title="Surplus Items"
                value={result.surplusItems}
                subtitle={result.surplusItems === 0 ? 'Perfect match!' : 'Minimized waste'}
                icon={TrendingDown}
                variant="orange"
            />

            <MetricCard
                title="Total Packs"
                value={result.totalPacks}
                subtitle="Optimized for shipping"
                icon={Package2}
                variant="green"
            />
        </div>
    );
} 