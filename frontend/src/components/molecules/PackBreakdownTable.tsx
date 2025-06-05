import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { formatNumber } from '@/lib/utils';
import type { CalculationResult } from '@/lib/types';

interface PackBreakdownTableProps {
    packBreakdown: CalculationResult['packBreakdown'];
}

export function PackBreakdownTable({ packBreakdown }: PackBreakdownTableProps) {
    return (
        <div className="border rounded-lg overflow-hidden">
            <Table>
                <TableHeader>
                    <TableRow className="bg-muted/50">
                        <TableHead className="font-semibold">Pack Size</TableHead>
                        <TableHead className="font-semibold text-right">Quantity Needed</TableHead>
                        <TableHead className="font-semibold text-right">Total Items</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {packBreakdown.map((pack) => (
                        <TableRow key={pack.packSize}>
                            <TableCell className="font-medium">
                                {formatNumber(pack.packSize)}
                            </TableCell>
                            <TableCell className="text-right font-semibold">
                                {formatNumber(pack.quantity)}
                            </TableCell>
                            <TableCell className="text-right">
                                {formatNumber(pack.packSize * pack.quantity)}
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </div>
    );
} 