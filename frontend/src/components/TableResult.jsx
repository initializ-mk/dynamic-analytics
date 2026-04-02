import { useState, useMemo } from 'react';

const INITIAL_ROWS = 20;

export default function TableResult({ data, columns }) {
  const [showAll, setShowAll] = useState(false);
  const [sortField, setSortField] = useState(null);
  const [sortDir, setSortDir] = useState('asc');

  // Derive columns from data if not provided
  const cols = useMemo(() => {
    if (columns && columns.length > 0) return columns;
    if (!data || data.length === 0) return [];
    return Object.keys(data[0]).map((key) => ({
      field: key,
      header: key.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase()),
    }));
  }, [columns, data]);

  const sorted = useMemo(() => {
    if (!sortField || !data) return data;
    return [...data].sort((a, b) => {
      const aVal = a[sortField];
      const bVal = b[sortField];
      if (aVal == null && bVal == null) return 0;
      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (typeof aVal === 'number' && typeof bVal === 'number') {
        return sortDir === 'asc' ? aVal - bVal : bVal - aVal;
      }
      return sortDir === 'asc'
        ? String(aVal).localeCompare(String(bVal))
        : String(bVal).localeCompare(String(aVal));
    });
  }, [data, sortField, sortDir]);

  if (!data || data.length === 0) return <p className="text-slate-500">No data to display</p>;

  const visible = showAll ? sorted : sorted.slice(0, INITIAL_ROWS);

  const handleSort = (field) => {
    if (sortField === field) {
      setSortDir((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      setSortField(field);
      setSortDir('asc');
    }
  };

  const formatCell = (value, format) => {
    if (value == null) return '—';
    const num = Number(value);
    if (format === 'number' && !isNaN(num)) return new Intl.NumberFormat('en-US').format(num);
    if (format === 'percent' && !isNaN(num)) return `${num.toFixed(1)}%`;
    if (format === 'currency' && !isNaN(num)) return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(num);
    if (format === 'duration' && !isNaN(num)) {
      if (num >= 60) return `${Math.floor(num / 60)}m ${Math.floor(num % 60)}s`;
      return `${Math.floor(num)}s`;
    }
    if (typeof value === 'object') return JSON.stringify(value);
    return String(value);
  };

  return (
    <div>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-slate-700">
              {cols.map((col) => (
                <th
                  key={col.field}
                  onClick={() => handleSort(col.field)}
                  className="px-4 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider cursor-pointer hover:text-slate-200 transition-colors select-none"
                >
                  {col.header}
                  {sortField === col.field && (
                    <span className="ml-1">{sortDir === 'asc' ? '\u2191' : '\u2193'}</span>
                  )}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {visible.map((row, i) => (
              <tr key={i} className={`border-b border-slate-700/50 ${i % 2 === 0 ? 'bg-slate-800/30' : ''}`}>
                {cols.map((col) => (
                  <td key={col.field} className="px-4 py-3 text-slate-300">
                    {formatCell(row[col.field], col.format)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {data.length > INITIAL_ROWS && (
        <button
          onClick={() => setShowAll(!showAll)}
          className="mt-3 text-sm text-emerald-400 hover:text-emerald-300 transition-colors"
        >
          {showAll ? `Show first ${INITIAL_ROWS}` : `Show all ${data.length} rows`}
        </button>
      )}
    </div>
  );
}
