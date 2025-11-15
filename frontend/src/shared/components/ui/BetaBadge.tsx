export function BetaBadge() {
  return (
    <div className="relative inline-flex">
      <span className="absolute -top-0.5 -right-0.5 flex h-2 w-2 z-10">
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
        <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
      </span>
      <span className="inline-flex items-center px-1.5 py-0.5 rounded-md text-[9px] font-bold bg-linear-to-r from-primary to-primary/60 text-primary-foreground shadow-sm border border-primary/20">
        BETA
      </span>
    </div>
  );
}

export function BetaBadgeLarge() {
  return (
    <div className="relative inline-flex">
      <span className="absolute -top-1 -right-1 flex h-4 w-4 z-10">
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
        <span className="relative inline-flex rounded-full h-4 w-4 bg-primary"></span>
      </span>
      <span className="inline-flex items-center px-4 py-2 rounded-full text-sm font-bold bg-linear-to-r from-primary via-primary/80 to-primary/60 text-primary-foreground shadow-xl border-2 border-primary/30 backdrop-blur-sm">
        <span className="mr-2">🚀</span>
        BETA VERSION
      </span>
    </div>
  );
}
