import { Calculator } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { LoadingButton } from '@/components/atoms/LoadingButton';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

interface ActionButtonsProps {
    onClear: () => void;
    onRecalculate: () => void;
    isRecalculateDisabled: boolean;
}

export function ActionButtons({
    onClear,
    onRecalculate,
    isRecalculateDisabled
}: ActionButtonsProps) {
    return (
        <div className="flex gap-2">
            <Button onClick={onClear} className="flex items-center gap-2">
                Clear Results
            </Button>
            <Tooltip>
                <TooltipTrigger asChild>
                    <LoadingButton
                        onClick={onRecalculate}
                        disabled={isRecalculateDisabled}
                        icon={Calculator}
                    >
                        Recalculate
                    </LoadingButton>
                </TooltipTrigger>
                <TooltipContent>
                    Recalculate with current pack sizes and order quantity
                </TooltipContent>
            </Tooltip>
        </div>
    );
} 