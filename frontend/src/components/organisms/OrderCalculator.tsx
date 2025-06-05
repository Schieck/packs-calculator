import { Calculator, CheckCircle, ChevronDown, ChevronUp } from 'lucide-react';
import { useState, useEffect } from 'react';

import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';

import { StatusAlert } from '@/components/atoms/StatusAlert';
import { OrderInputForm } from '@/components/molecules/OrderInputForm';
import { PackBreakdownTable } from '@/components/molecules/PackBreakdownTable';
import { SummaryMetrics } from '@/components/molecules/SummaryMetrics';
import { PerformanceInfo } from '@/components/molecules/PerformanceInfo';
import { ActionButtons } from '@/components/molecules/ActionButtons';

import { useCalculation } from '@/lib/hooks';
import { type OrderCalculationFormData } from '@/lib/schemas';

interface OrderCalculatorProps {
    packSizes: number[];
}

export function OrderCalculator({ packSizes }: OrderCalculatorProps) {
    const { result, isLoading, error, calculate, clearResult } = useCalculation();
    const [isInputCollapsed, setIsInputCollapsed] = useState(false);
    const [lastOrderData, setLastOrderData] = useState<OrderCalculationFormData | null>(null);

    useEffect(() => {
        if (result) {
            setIsInputCollapsed(true);
        }
    }, [result]);

    const onSubmit = async (data: OrderCalculationFormData) => {
        try {
            const quantity = parseInt(data.orderQuantity, 10);
            setLastOrderData(data); // Store the form data for recalculation
            await calculate(quantity, packSizes);
        } catch (err) {
            console.error('Calculation failed:', err);
        }
    };

    const handleClearResult = () => {
        clearResult();
        setIsInputCollapsed(false);
        setLastOrderData(null); // Clear the stored order data
    };

    const toggleInputCollapse = () => {
        setIsInputCollapsed(!isInputCollapsed);
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
        } else if (e.key === 'Escape') {
            clearResult();
        }
    };

    const isCalculateDisabled =
        isLoading ||
        packSizes.length === 0;

    return (
        <Card className="w-full">
            {!isInputCollapsed && (
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Calculator className="h-5 w-5 text-blue-600" />
                        Order Calculator
                    </CardTitle>
                </CardHeader>
            )}
            <CardContent className="space-y-6">
                {/* Error Display */}
                {error && (
                    <StatusAlert variant="destructive">
                        {error}
                    </StatusAlert>
                )}

                {/* Pack Size Warning */}
                {packSizes.length === 0 && (
                    <StatusAlert>
                        Please add at least one pack size before calculating orders.
                    </StatusAlert>
                )}

                {/* Order Input Form */}
                {!isInputCollapsed && (
                    <div className="space-y-4">
                        <OrderInputForm
                            onSubmit={onSubmit}
                            isLoading={isLoading}
                            isDisabled={isCalculateDisabled}
                            onKeyDown={handleKeyDown}
                        />
                        <Separator />
                    </div>
                )}

                {/* Results Section */}
                {result && (
                    <div className="space-y-6">
                        {/* Success Message */}
                        <StatusAlert
                            className="border-green-200 bg-green-50 cursor-pointer hover:bg-green-100 transition-colors"
                            onClick={toggleInputCollapse}
                            icon={CheckCircle}
                        >
                            <div className="flex items-center justify-between w-full text-green-800">
                                <span>Optimal pack combination calculated successfully!</span>
                                {isInputCollapsed ? (
                                    <ChevronDown className="h-4 w-4 text-green-600 ml-2" />
                                ) : (
                                    <ChevronUp className="h-4 w-4 text-green-600 ml-2" />
                                )}
                            </div>
                        </StatusAlert>

                        {/* Pack Breakdown Table */}
                        <div className="space-y-3">
                            <PackBreakdownTable packBreakdown={result.packBreakdown} />
                        </div>

                        {/* Summary Metrics */}
                        <SummaryMetrics result={result} />

                        {/* Performance Metrics */}
                        <PerformanceInfo
                            calculationTime={result.calculationTime}
                            orderQuantity={lastOrderData?.orderQuantity || "0"}
                            packSizesCount={packSizes.length}
                        />

                        {/* Action Buttons */}
                        <ActionButtons
                            onClear={handleClearResult}
                            onRecalculate={() => lastOrderData && onSubmit(lastOrderData)}
                            isRecalculateDisabled={isCalculateDisabled || !lastOrderData}
                        />
                    </div>
                )}
            </CardContent>
        </Card>
    );
} 