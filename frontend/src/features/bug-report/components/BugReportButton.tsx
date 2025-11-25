"use client";

import { useState } from "react";
import { Bug } from "lucide-react";
import { BugReportDialog } from "./BugReportDialog";

interface BugReportButtonProps {
  variant?: "full" | "icon" | "header";
  className?: string;
}

export function BugReportButton({ variant = "full", className = "" }: BugReportButtonProps) {
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  if (variant === "header") {
    return (
      <>
        <button
          onClick={() => setIsDialogOpen(true)}
          className={`flex items-center gap-3 p-2 md:py-1.5 rounded-lg bg-muted hover:bg-muted/80 transition-colors text-sm text-foreground ${className}`}
          title="Report a bug"
        >
          <Bug className="w-5 h-5 md:w-4 md:h-4" />
          <span className="hidden sm:inline">Report Bug</span>
        </button>
        <BugReportDialog isOpen={isDialogOpen} onClose={() => setIsDialogOpen(false)} />
      </>
    );
  }

  if (variant === "icon") {
    return (
      <>
        <button
          onClick={() => setIsDialogOpen(true)}
          className={`p-3 rounded-lg transition-colors relative group cursor-pointer hover:bg-secondary/50 text-muted-foreground hover:text-foreground ${className}`}
          title="Report a Bug"
        >
          <Bug className="w-5 h-5" />
          <div className="absolute left-full ml-2 top-1/2 -translate-y-1/2 px-2 py-1 bg-popover text-popover-foreground text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none whitespace-nowrap z-50">
            Report a Bug
          </div>
        </button>
        <BugReportDialog isOpen={isDialogOpen} onClose={() => setIsDialogOpen(false)} />
      </>
    );
  }

  return (
    <>
      <button
        onClick={() => setIsDialogOpen(true)}
        className={`w-full px-4 py-3 rounded-lg flex items-center gap-3 transition-colors cursor-pointer hover:bg-secondary/50 text-muted-foreground hover:text-foreground ${className}`}
      >
        <Bug className="w-5 h-5" />
        <span className="text-sm font-semibold">Report a Bug</span>
      </button>
      <BugReportDialog isOpen={isDialogOpen} onClose={() => setIsDialogOpen(false)} />
    </>
  );
}
