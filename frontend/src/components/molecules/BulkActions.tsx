import { RotateCcw, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';

interface BulkActionsProps {
    onResetToDefaults: () => void;
    onClearAll: () => void;
}

export function BulkActions({ onResetToDefaults, onClearAll }: BulkActionsProps) {
    return (
        <div className="space-y-3">
            <Label className="text-sm font-medium">Bulk Actions</Label>
            <div className="flex flex-wrap gap-2">
                <Tooltip>
                    <TooltipTrigger asChild>
                        <Button
                            size="sm"
                            onClick={onResetToDefaults}
                            className="flex items-center gap-2"
                        >
                            <RotateCcw className="h-4 w-4" />
                            Reset to Defaults
                        </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                        Reset to standard pack sizes: 250, 500, 1000, 2000, 5000
                    </TooltipContent>
                </Tooltip>
                <Tooltip>
                    <TooltipTrigger asChild>
                        <Button
                            size="sm"
                            onClick={onClearAll}
                            className="flex items-center gap-2 text-destructive hover:bg-destructive hover:text-destructive-foreground"
                        >
                            <Trash2 className="h-4 w-4" />
                            Clear All
                        </Button>
                    </TooltipTrigger>
                    <TooltipContent>
                        Remove all pack sizes (you'll need to add new ones)
                    </TooltipContent>
                </Tooltip>
            </div>
        </div>
    );
} 