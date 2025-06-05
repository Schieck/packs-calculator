import { Clock } from 'lucide-react';
import { formatTime, formatNumber } from '@/lib/utils';

interface PerformanceInfoProps {
    calculationTime: number;
    orderQuantity: string;
    packSizesCount: number;
}

export function PerformanceInfo({
    calculationTime,
    orderQuantity,
    packSizesCount
}: PerformanceInfoProps) {
    return (
        <div className="bg-gray-50 border border-gray-200 rounded-lg p-4">
            <div className="flex items-center gap-2 mb-2">
                <Clock className="h-4 w-4 text-gray-600" />
                <span className="text-sm font-medium text-gray-900">Performance</span>
            </div>
            <div className="text-sm text-gray-700">
                Calculation completed in{' '}
                <span className="font-semibold">{formatTime(calculationTime)}</span>
            </div>
            <div className="text-xs text-gray-500 mt-1">
                Processed {formatNumber(parseInt(orderQuantity || '0'))} items with {packSizesCount} pack sizes
            </div>
        </div>
    );
} 