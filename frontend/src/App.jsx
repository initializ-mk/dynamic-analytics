import { useState, useCallback } from 'react';
import QueryInput from './components/QueryInput';
import ResultRenderer from './components/ResultRenderer';
import LoadingState from './components/LoadingState';
import QueryHistory from './components/QueryHistory';
import { executeQuery } from './api/query';

export default function App() {
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [history, setHistory] = useState([]);

  const handleQuery = useCallback(async (queryText) => {
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const data = await executeQuery(queryText);
      setResult(data);
      setHistory((prev) => {
        const filtered = prev.filter((q) => q !== queryText);
        return [queryText, ...filtered].slice(0, 5);
      });
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <div className="min-h-screen bg-slate-900">
      {/* Header */}
      <header className="border-b border-slate-700 bg-slate-900/80 backdrop-blur-sm sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 rounded-lg bg-emerald-500/20 flex items-center justify-center">
              <svg className="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 3v11.25A2.25 2.25 0 0 0 6 16.5h2.25M3.75 3h-1.5m1.5 0h16.5m0 0h1.5m-1.5 0v11.25A2.25 2.25 0 0 1 18 16.5h-2.25m-7.5 0h7.5m-7.5 0-1 3m8.5-3 1 3m0 0 .5 1.5m-.5-1.5h-9.5m0 0-.5 1.5m.75-9 3-3 2.148 2.148A12.061 12.061 0 0 1 16.5 7.605" />
              </svg>
            </div>
            <div>
              <h1 className="text-lg font-semibold text-slate-100">Recruitment Analytics</h1>
              <p className="text-xs text-slate-400">AI-powered insights from your hiring pipeline</p>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <QueryInput onSubmit={handleQuery} disabled={loading} />

        {history.length > 0 && (
          <QueryHistory queries={history} onSelect={handleQuery} />
        )}

        {loading && <LoadingState />}

        {error && (
          <div className="mt-6 p-4 rounded-xl bg-red-500/10 border border-red-500/20">
            <p className="text-red-400 text-sm font-medium">Error</p>
            <p className="text-red-300 text-sm mt-1">{error}</p>
          </div>
        )}

        {result && !loading && <ResultRenderer result={result} />}
      </main>
    </div>
  );
}
