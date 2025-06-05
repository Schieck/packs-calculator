import React from 'react';
import { AlertTriangle, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';

interface ErrorBoundaryProps {
    children: React.ReactNode;
    fallback?: React.ComponentType<{ error: Error; resetError: () => void }>;
}

interface ErrorBoundaryState {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends React.Component<
    ErrorBoundaryProps,
    ErrorBoundaryState
> {
    constructor(props: ErrorBoundaryProps) {
        super(props);
        this.state = { hasError: false, error: null };
    }

    static getDerivedStateFromError(error: Error): ErrorBoundaryState {
        return {
            hasError: true,
            error,
        };
    }

    componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
        console.error('Error caught by boundary:', error, errorInfo);
    }

    resetError = () => {
        this.setState({ hasError: false, error: null });
    };

    render() {
        if (this.state.hasError) {
            if (this.props.fallback) {
                const FallbackComponent = this.props.fallback;
                return (
                    <FallbackComponent
                        error={this.state.error!}
                        resetError={this.resetError}
                    />
                );
            }

            return (
                <DefaultErrorFallback
                    error={this.state.error!}
                    resetError={this.resetError}
                />
            );
        }

        return this.props.children;
    }
}

interface ErrorFallbackProps {
    error: Error;
    resetError: () => void;
}

function DefaultErrorFallback({ error, resetError }: ErrorFallbackProps) {
    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 flex items-center justify-center p-4">
            <Card className="w-full max-w-lg">
                <CardHeader>
                    <CardTitle className="flex items-center gap-2 text-destructive">
                        <AlertTriangle className="h-5 w-5" />
                        Something Went Wrong
                    </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                    <Alert variant="destructive">
                        <AlertDescription>
                            An unexpected error occurred while running the pack calculator.
                            This may be due to a network issue or a temporary problem.
                        </AlertDescription>
                    </Alert>

                    <div className="bg-muted p-4 rounded-lg">
                        <h4 className="text-sm font-medium mb-2">Error Details:</h4>
                        <code className="text-xs text-muted-foreground break-all">
                            {error.message}
                        </code>
                    </div>

                    <div className="flex gap-2">
                        <Button
                            onClick={resetError}
                            className="flex items-center gap-2"
                        >
                            <RefreshCw className="h-4 w-4" />
                            Try Again
                        </Button>
                        <Button
                            onClick={() => window.location.reload()}
                            className="flex items-center gap-2"
                        >
                            Reload Page
                        </Button>
                    </div>

                    <div className="text-xs text-muted-foreground">
                        If this problem persists, please check your network connection
                        or contact support.
                    </div>
                </CardContent>
            </Card>
        </div>
    );
} 