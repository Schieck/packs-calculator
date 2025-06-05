import { Label } from '@/components/ui/label';
import { PackSizeBadge } from '@/components/atoms/PackSizeBadge';

interface PackSizeListProps {
    packSizes: number[];
    onRemove: (size: number) => void;
}

export function PackSizeList({ packSizes, onRemove }: PackSizeListProps) {
    return (
        <div className="space-y-3">
            <Label className="text-sm font-medium">Current Pack Sizes</Label>
            <div className="flex flex-wrap gap-2 min-h-[2rem]">
                {packSizes.length === 0 ? (
                    <div className="text-sm text-muted-foreground italic">
                        No pack sizes configured. Add some pack sizes to get started.
                    </div>
                ) : (
                    packSizes.map((size) => (
                        <PackSizeBadge
                            key={size}
                            size={size}
                            onRemove={onRemove}
                        />
                    ))
                )}
            </div>
        </div>
    );
} 