import { Sparkles } from "lucide-react";

export function ChatEmptyState() {
  return (
    <div className="flex flex-col items-center justify-center h-full space-y-4 text-center pt-20">
      <Sparkles className="w-16 h-16 text-primary/50" />
      <h2 className="text-2xl font-bold">What are you looking for?</h2>
      <p className="text-muted-foreground max-w-md">
        Tell me what product you need and I'll help you find the best options
      </p>
    </div>
  );
}
