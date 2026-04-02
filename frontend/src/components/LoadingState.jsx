export default function LoadingState() {
  return (
    <div className="mt-6 bg-slate-800 border border-slate-700 rounded-xl p-6 animate-pulse">
      <div className="h-5 bg-slate-700 rounded w-1/3 mb-3" />
      <div className="h-3 bg-slate-700 rounded w-2/3 mb-6" />
      <div className="space-y-3">
        <div className="h-48 bg-slate-700/50 rounded-lg" />
        <div className="flex gap-4">
          <div className="h-3 bg-slate-700 rounded w-16" />
          <div className="h-3 bg-slate-700 rounded w-12" />
        </div>
      </div>
    </div>
  );
}
