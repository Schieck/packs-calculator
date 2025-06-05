import { X } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { formatNumber } from '@/lib/utils';

interface PackSizeBadgeProps {
    size: number;
    onRemove: (size: number) => void;
}

export function PackSizeBadge({ size, onRemove }: PackSizeBadgeProps) {
    return (
        <Badge
            variant="secondary"
            className="flex items-center gap-1 px-3 py-1 text-sm"
        >
            {formatNumber(size)}
            <Tooltip>
                <TooltipTrigger asChild>
                    <Button
                        onClick={() => onRemove(size)}
                        className="ml-1 hover:bg-destructive hover:text-destructive-foreground rounded-full p-0.5 transition-colors"
                        aria-label={`Remove pack size ${size}`}
                    >
                        <X className="h-3 w-3" />
                    </Button>
                </TooltipTrigger>
                <TooltipContent>
                    Remove pack size {formatNumber(size)}
                </TooltipContent>
            </Tooltip>
        </Badge>
    );
} 