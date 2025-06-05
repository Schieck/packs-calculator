import { TrendingUp } from 'lucide-react';
import { TooltipProvider } from '@/components/TooltipProvider';
import { PackSizeManager, OrderCalculator } from '@/components/organisms';
import { usePackSizes } from '@/lib/store';
import './App.css';

function App() {
  const packSizes = usePackSizes();

  return (
    <TooltipProvider>
      <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
        {/* Main Content */}
        <main className="container mx-auto px-4 py-8">
          <div className="max-w-6xl mx-auto space-y-8">
            {/* Hero Section */}
            <div className="text-center space-y-4 mb-12">
              <div className="flex items-center justify-center gap-2 text-blue-600 mb-4">
                <TrendingUp className="h-6 w-6" />
                <span className="text-lg font-semibold">
                  Pack Optimization Calculator
                </span>
              </div>
            </div>

            {/* Application Cards */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Pack Size Management */}
              <div className="lg:col-span-1">
                <PackSizeManager />
              </div>

              {/* Order Calculator */}
              <div className="lg:col-span-1">
                <OrderCalculator packSizes={packSizes} />
              </div>
            </div>
          </div>
        </main>

        {/* Footer */}
        <footer className="bg-white/80 backdrop-blur-sm border-t border-slate-200 mt-16">
          <div className="container mx-auto px-4 py-8">
            <div className="text-center">
              <p className="text-sm text-slate-600">
                Built with 💜 and modern engineering practices by @Schieck.
              </p>
            </div>
          </div>
        </footer>
      </div>
    </TooltipProvider>
  );
}

export default App;
