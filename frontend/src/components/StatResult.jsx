export default function StatResult({ data, config }) {
  if (!data || data.length === 0) return <p className="text-slate-500">No data</p>;

  const valueField = config?.value_field || Object.keys(data[0]).find((k) => k !== '_id') || Object.keys(data[0])[0];
  const rawValue = data[0][valueField];
  const format = config?.format || 'number';
  const label = config?.label || valueField;

  const formatValue = (val, fmt) => {
    if (val == null) return '—';
    const num = Number(val);
    if (isNaN(num)) return String(val);

    switch (fmt) {
      case 'percent':
        return `${num.toFixed(1)}%`;
      case 'duration':
        if (num >= 3600) return `${Math.floor(num / 3600)}h ${Math.floor((num % 3600) / 60)}m`;
        if (num >= 60) return `${Math.floor(num / 60)}m ${Math.floor(num % 60)}s`;
        return `${Math.floor(num)}s`;
      case 'currency':
        return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(num);
      default:
        return num >= 1000 ? new Intl.NumberFormat('en-US').format(Math.round(num)) : String(Math.round(num * 100) / 100);
    }
  };

  return (
    <div className="flex flex-col items-center py-6">
      <span className="text-5xl font-bold text-emerald-400">
        {formatValue(rawValue, format)}
      </span>
      <span className="text-sm text-slate-400 mt-2">{label}</span>
    </div>
  );
}
