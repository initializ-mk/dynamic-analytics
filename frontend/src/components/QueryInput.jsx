import { useState } from 'react';

const SUGGESTED_QUERIES = [
  'Total candidates',
  'Candidates by region',
  'Candidates by status',
  'Hiring funnel breakdown',
  'Top recruiters by screenings',
  'Average call duration by status',
  'Candidates by industry',
  'Weekly candidate trend',
  'Senior engineers by region',
  'Conversion rate from screened to hired',
];

export default function QueryInput({ onSubmit, disabled }) {
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (query.trim() && !disabled) {
      onSubmit(query.trim());
    }
  };

  const handlePillClick = (q) => {
    if (!disabled) {
      setQuery(q);
      onSubmit(q);
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="relative">
        <div className="flex items-center gap-3 bg-slate-800 border border-slate-700 rounded-xl px-4 py-3 focus-within:border-emerald-500/50 focus-within:ring-1 focus-within:ring-emerald-500/25 transition-all">
          <svg className="w-5 h-5 text-slate-400 flex-shrink-0" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
          </svg>
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Ask anything about your recruitment data..."
            className="flex-1 bg-transparent text-slate-100 placeholder-slate-500 outline-none text-sm"
            disabled={disabled}
          />
          <button
            type="submit"
            disabled={disabled || !query.trim()}
            className="px-4 py-1.5 bg-emerald-500 text-white text-sm font-medium rounded-lg hover:bg-emerald-400 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {disabled ? 'Analyzing...' : 'Ask'}
          </button>
        </div>
      </form>

      <div className="mt-4 flex flex-wrap gap-2">
        {SUGGESTED_QUERIES.map((q) => (
          <button
            key={q}
            onClick={() => handlePillClick(q)}
            disabled={disabled}
            className="px-3 py-1.5 text-xs bg-slate-800 text-slate-300 border border-slate-700 rounded-full hover:bg-slate-700 hover:text-slate-100 disabled:opacity-40 transition-colors"
          >
            {q}
          </button>
        ))}
      </div>
    </div>
  );
}
