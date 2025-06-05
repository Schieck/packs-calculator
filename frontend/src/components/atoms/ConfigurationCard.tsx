import { Star, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { PackConfiguration } from '@/lib/types';

interface ConfigurationCardProps {
    configuration: PackConfiguration;
    onSelect: (config: PackConfiguration) => void;
    onDelete: (id: number) => void;
    onSetDefault: (id: number) => void;
    isSelected?: boolean;
    isLoading?: boolean;
}

export function ConfigurationCard({
    configuration,
    onSelect,
    onDelete,
    onSetDefault,
    isSelected = false,
    isLoading = false,
}: ConfigurationCardProps) {
    return (
        <Card
            className={`cursor-pointer transition-all hover:shadow-md ${isSelected ? 'ring-2 ring-blue-500 bg-blue-50' : ''
                }`}
            onClick={() => onSelect(configuration)}
        >
            <CardContent className="p-3">
                <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                        <h4 className="font-medium text-sm">{configuration.name}</h4>
                        {configuration.is_default && (
                            <Badge variant="secondary" className="text-xs">
                                <Star className="h-3 w-3 mr-1" />
                                Default
                            </Badge>
                        )}
                    </div>
                    <div className="flex gap-1" onClick={(e) => e.stopPropagation()}>
                        {!configuration.is_default && (
                            <Button
                                size="sm"
                                onClick={() => onSetDefault(configuration.id)}
                                disabled={isLoading}
                                className="h-6 w-6 p-0"
                                title="Set as default"
                            >
                                <Star className="h-3 w-3" />
                            </Button>
                        )}
                        <Button
                            size="sm"
                            onClick={() => onDelete(configuration.id)}
                            disabled={isLoading || configuration.is_default}
                            className="h-6 w-6 p-0 text-red-500 hover:text-red-700"
                            title="Delete configuration"
                        >
                            <Trash2 className="h-3 w-3" />
                        </Button>
                    </div>
                </div>
                <div className="flex flex-wrap gap-1">
                    {configuration.pack_sizes.map((size) => (
                        <Badge key={size} variant="outline" className="text-xs">
                            {size}
                        </Badge>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
} 