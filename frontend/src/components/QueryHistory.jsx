export default function QueryHistory({ queries, onSelect }) {
  return (
    <div className="mt-4 flex items-center gap-2 overflow-x-auto pb-2">
      <span className="text-xs text-slate-500 flex-shrink-0">Recent:</span>
      {queries.map((q) => (
        <button
          key={q}
          onClick={() => onSelect(q)}
          className="flex-shrink-0 px-3 py-1 text-xs bg-slate-800/50 text-slate-400 border border-slate-700/50 rounded-full hover:bg-slate-700 hover:text-slate-200 transition-colors truncate max-w-[200px]"
          title={q}
        >
          {q}
        </button>
      ))}
    </div>
  );
}
