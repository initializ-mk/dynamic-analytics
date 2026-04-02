import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  LineChart, Line, PieChart, Pie, Cell, Legend,
} from 'recharts';

const COLORS = [
  '#10b981', '#3b82f6', '#f59e0b', '#ef4444',
  '#8b5cf6', '#ec4899', '#06b6d4', '#f97316',
];

const tooltipStyle = {
  contentStyle: { backgroundColor: '#1e293b', border: '1px solid #334155', borderRadius: '8px' },
  labelStyle: { color: '#94a3b8' },
  itemStyle: { color: '#e2e8f0' },
};

export default function ChartResult({ data, config, chartType }) {
  if (!data || data.length === 0) return <p className="text-slate-500">No data to display</p>;

  const xField = config?.x_field || '_id';
  const yField = config?.y_field || Object.keys(data[0]).find((k) => k !== '_id' && k !== xField) || 'count';
  const xLabel = config?.x_label || xField;
  const yLabel = config?.y_label || yField;

  // Normalize data: ensure _id values are strings for chart rendering
  const chartData = data.map((d) => ({
    ...d,
    [xField]: d[xField] != null ? String(d[xField]) : 'N/A',
  }));

  if (chartType === 'pie_chart') {
    return (
      <ResponsiveContainer width="100%" height={400}>
        <PieChart>
          <Pie
            data={chartData}
            dataKey={yField}
            nameKey={xField}
            cx="50%"
            cy="50%"
            outerRadius={140}
            label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
            labelLine={true}
          >
            {chartData.map((_, i) => (
              <Cell key={i} fill={COLORS[i % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip {...tooltipStyle} />
          <Legend wrapperStyle={{ color: '#94a3b8' }} />
        </PieChart>
      </ResponsiveContainer>
    );
  }

  if (chartType === 'line_chart') {
    return (
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis dataKey={xField} stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} label={{ value: xLabel, position: 'insideBottom', offset: -5, fill: '#94a3b8' }} />
          <YAxis stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} label={{ value: yLabel, angle: -90, position: 'insideLeft', fill: '#94a3b8' }} />
          <Tooltip {...tooltipStyle} />
          <Line type="monotone" dataKey={yField} stroke="#10b981" strokeWidth={2} dot={{ fill: '#10b981', r: 4 }} activeDot={{ r: 6 }} />
        </LineChart>
      </ResponsiveContainer>
    );
  }

  if (chartType === 'funnel') {
    // Render funnel as horizontal decreasing bar chart sorted by count desc
    const sorted = [...chartData].sort((a, b) => (b[yField] || 0) - (a[yField] || 0));
    return (
      <ResponsiveContainer width="100%" height={Math.max(300, sorted.length * 50)}>
        <BarChart data={sorted} layout="vertical" margin={{ top: 5, right: 30, left: 80, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
          <XAxis type="number" stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} />
          <YAxis type="category" dataKey={xField} stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} width={100} />
          <Tooltip {...tooltipStyle} />
          <Bar dataKey={yField} radius={[0, 4, 4, 0]}>
            {sorted.map((_, i) => (
              <Cell key={i} fill={COLORS[i % COLORS.length]} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    );
  }

  // Default: bar_chart
  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={chartData} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
        <XAxis dataKey={xField} stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} label={{ value: xLabel, position: 'insideBottom', offset: -5, fill: '#94a3b8' }} />
        <YAxis stroke="#64748b" tick={{ fill: '#94a3b8', fontSize: 12 }} label={{ value: yLabel, angle: -90, position: 'insideLeft', fill: '#94a3b8' }} />
        <Tooltip {...tooltipStyle} />
        <Bar dataKey={yField} fill="#10b981" radius={[4, 4, 0, 0]}>
          {chartData.map((_, i) => (
            <Cell key={i} fill={COLORS[i % COLORS.length]} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
