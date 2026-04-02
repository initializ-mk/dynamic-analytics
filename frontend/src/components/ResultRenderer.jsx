import ChartResult from './ChartResult';
import TableResult from './TableResult';
import StatResult from './StatResult';

export default function ResultRenderer({ result }) {
  const { ui_type, data, title, summary, chart_config, columns, stat_config, generated_pipeline, meta } = result;

  return (
    <div className="mt-6 space-y-4">
      {/* Title and summary */}
      <div className="bg-slate-800 border border-slate-700 rounded-xl p-6">
        {title && <h2 className="text-lg font-semibold text-slate-100">{title}</h2>}
        {summary && <p className="text-sm text-slate-400 mt-1">{summary}</p>}

        {/* Result visualization */}
        <div className="mt-5">
          {ui_type === 'stat' && (
            <StatResult data={data} config={stat_config} />
          )}
          {(ui_type === 'bar_chart' || ui_type === 'line_chart' || ui_type === 'pie_chart' || ui_type === 'funnel') && (
            <ChartResult data={data} config={chart_config} chartType={ui_type} />
          )}
          {ui_type === 'table' && (
            <TableResult data={data} columns={columns} />
          )}
        </div>

        {/* Meta info */}
        {meta && (
          <div className="mt-4 pt-4 border-t border-slate-700 flex items-center gap-4 text-xs text-slate-500">
            <span>{meta.result_count} result{meta.result_count !== 1 ? 's' : ''}</span>
            <span>{meta.execution_time_ms}ms</span>
          </div>
        )}
      </div>

      {/* Pipeline debug */}
      {generated_pipeline && (
        <details className="bg-slate-800 border border-slate-700 rounded-xl p-4">
          <summary className="text-xs text-slate-500 cursor-pointer hover:text-slate-300 transition-colors">
            View generated pipeline
          </summary>
          <pre className="mt-3 text-xs text-slate-400 overflow-x-auto">
            {JSON.stringify(generated_pipeline, null, 2)}
          </pre>
        </details>
      )}
    </div>
  );
}
