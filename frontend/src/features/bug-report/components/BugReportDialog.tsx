"use client";

import { useState, useEffect, useRef } from "react";
import { createPortal } from "react-dom";
import { X, Bug, Send, Loader2, Upload } from "lucide-react";
import { useAuthStore } from "@/shared/lib/auth-store";
import { useChatStore } from "@/shared/lib/store";

interface BugReportDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export function BugReportDialog({ isOpen, onClose }: BugReportDialogProps) {
  const [description, setDescription] = useState("");
  const [steps, setSteps] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitStatus, setSubmitStatus] = useState<"idle" | "success" | "error">("idle");
  const [mounted, setMounted] = useState(false);
  const [attachments, setAttachments] = useState<File[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const user = useAuthStore((state) => state.user);
  const sessionId = useChatStore((state) => state.sessionId);
  const messages = useChatStore((state) => state.messages);

  // Track when component is mounted (client-side only)
  useEffect(() => {
    setMounted(true);
  }, []);

  // Block body scroll when modal is open and prevent layout shift
  useEffect(() => {
    if (isOpen) {
      // Get scrollbar width before hiding
      const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth;

      // Block scroll and compensate for scrollbar width
      document.body.style.overflow = "hidden";
      document.body.style.paddingRight = `${scrollbarWidth}px`;
    } else {
      // Restore scroll
      document.body.style.overflow = "unset";
      document.body.style.paddingRight = "0px";
    }

    // Cleanup on unmount
    return () => {
      document.body.style.overflow = "unset";
      document.body.style.paddingRight = "0px";
    };
  }, [isOpen]);

  // Close on Escape key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen && !isSubmitting) {
        onClose();
      }
    };

    document.addEventListener("keydown", handleEscape);
    return () => document.removeEventListener("keydown", handleEscape);
  }, [isOpen, isSubmitting, onClose]);

  // Handle paste event for screenshots (Ctrl+V)
  useEffect(() => {
    const handlePaste = (e: ClipboardEvent) => {
      if (!isOpen) return;

      const items = e.clipboardData?.items;
      if (!items) return;

      for (let i = 0; i < items.length; i++) {
        const item = items[i];
        if (item.type.indexOf("image") !== -1) {
          const file = item.getAsFile();
          if (file) {
            handleFileAdd(file);
          }
        }
      }
    };

    document.addEventListener("paste", handlePaste);
    return () => document.removeEventListener("paste", handlePaste);
  }, [isOpen]);

  // File handling functions
  const handleFileAdd = (file: File) => {
    // Check if file is an image
    if (!file.type.startsWith("image/")) {
      alert("Only image files are allowed");
      return;
    }

    // Check file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert("File size must be less than 5MB");
      return;
    }

    // Check if already added
    if (attachments.some((f) => f.name === file.name && f.size === file.size)) {
      return;
    }

    // Add to attachments (max 3 files)
    if (attachments.length >= 3) {
      alert("Maximum 3 files allowed");
      return;
    }

    setAttachments((prev) => [...prev, file]);
  };

  const handleFileRemove = (index: number) => {
    setAttachments((prev) => prev.filter((_, i) => i !== index));
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    files.forEach((file) => handleFileAdd(file));
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    files.forEach((file) => handleFileAdd(file));
    // Reset input value to allow selecting the same file again
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSubmitStatus("idle");

    try {
      // Collect context
      const context = {
        user_id: user?.id || "anonymous",
        user_email: user?.email || "not_logged_in",
        session_id: sessionId || "no_session",
        url: window.location.href,
        user_agent: navigator.userAgent,
        screen_resolution: `${window.screen.width}x${window.screen.height}`,
        viewport_size: `${window.innerWidth}x${window.innerHeight}`,
        timestamp: new Date().toISOString(),
        last_messages: messages.slice(-5).map((m) => ({
          role: m.role,
          content: m.content?.substring(0, 100) + (m.content && m.content.length > 100 ? "..." : ""),
        })),
      };

      // Create FormData to send files
      const formData = new FormData();
      formData.append("description", description);
      formData.append("steps_to_reproduce", steps);
      formData.append("context", JSON.stringify(context));

      // Add attachments
      attachments.forEach((file, index) => {
        formData.append(`attachments`, file);
      });

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/bug-report`, {
        method: "POST",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Failed to submit bug report");
      }

      setSubmitStatus("success");
      setTimeout(() => {
        onClose();
        setDescription("");
        setSteps("");
        setAttachments([]);
        setSubmitStatus("idle");
      }, 2000);
    } catch (error) {
      console.error("Error submitting bug report:", error);
      setSubmitStatus("error");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!isOpen || !mounted) return null;

  const modalContent = (
    <div
      className="fixed inset-0 bg-black/50 backdrop-blur-sm z-[9999] flex items-center justify-center p-4 animate-fade-in-up"
      onClick={(e) => {
        // Close only if clicking on the backdrop, not the modal content
        if (e.target === e.currentTarget && !isSubmitting) {
          onClose();
        }
      }}
    >
      <div
        className="bg-background border border-border rounded-xl shadow-2xl w-full max-w-2xl max-h-[85vh] overflow-hidden flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="px-4 py-3 border-b border-border flex items-center justify-between bg-muted/30">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-destructive/10 flex items-center justify-center">
              <Bug className="w-4 h-4 text-destructive" />
            </div>
            <div>
              <h2 className="text-base font-semibold">Report a Bug</h2>
              <p className="text-[11px] text-muted-foreground">Help us improve</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-md hover:bg-secondary/80 transition-colors"
            disabled={isSubmitting}
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto p-4 space-y-3">
          {/* Description */}
          <div className="space-y-1.5">
            <label htmlFor="description" className="text-xs font-medium text-foreground/90">
              What happened? <span className="text-destructive">*</span>
            </label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Describe the bug you encountered..."
              required
              disabled={isSubmitting}
              className="w-full min-h-[70px] px-2.5 py-2 rounded-md border border-border bg-background text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 disabled:opacity-50 resize-none"
            />
          </div>

          {/* Steps to reproduce */}
          <div className="space-y-1.5">
            <label htmlFor="steps" className="text-xs font-medium text-foreground/90">
              Steps to reproduce (optional)
            </label>
            <textarea
              id="steps"
              value={steps}
              onChange={(e) => setSteps(e.target.value)}
              placeholder="1. Go to...&#10;2. Click on...&#10;3. See error..."
              disabled={isSubmitting}
              className="w-full min-h-[60px] px-2.5 py-2 rounded-md border border-border bg-background text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50 disabled:opacity-50 resize-none"
            />
          </div>

          {/* File attachments */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-foreground/90">
              Screenshots (optional)
            </label>

            {/* Drag and drop area */}
            <div
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              className={`relative border border-dashed rounded-md p-3 text-center transition-colors ${
                isDragging
                  ? "border-primary bg-primary/5"
                  : "border-border bg-secondary/10 hover:bg-secondary/20"
              }`}
            >
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                multiple
                onChange={handleFileInputChange}
                className="hidden"
                disabled={isSubmitting}
              />

              <div className="flex flex-col items-center gap-1.5">
                <Upload className="w-6 h-6 text-muted-foreground" />
                <div className="text-xs">
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    disabled={isSubmitting}
                    className="text-primary hover:underline font-medium"
                  >
                    Click to upload
                  </button>
                  <span className="text-muted-foreground"> or drag & drop</span>
                </div>
                <p className="text-[11px] text-muted-foreground">
                  Up to 3 images, 5MB each • Press <kbd className="px-1 py-0.5 bg-muted rounded text-[10px]">Ctrl+V</kbd> to paste
                </p>
              </div>
            </div>

            {/* Preview attachments */}
            {attachments.length > 0 && (
              <div className="grid grid-cols-3 gap-2">
                {attachments.map((file, index) => (
                  <div
                    key={`${file.name}-${index}`}
                    className="relative group aspect-square rounded-md overflow-hidden border border-border bg-secondary/20"
                  >
                    <img
                      src={URL.createObjectURL(file)}
                      alt={file.name}
                      className="w-full h-full object-cover"
                    />
                    <div className="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <button
                        type="button"
                        onClick={() => handleFileRemove(index)}
                        disabled={isSubmitting}
                        className="p-1.5 bg-destructive text-destructive-foreground rounded-full hover:bg-destructive/90 transition-colors"
                      >
                        <X className="w-3.5 h-3.5" />
                      </button>
                    </div>
                    <div className="absolute bottom-0 left-0 right-0 bg-black/70 text-white text-[10px] px-1 py-0.5 truncate">
                      {file.name}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Auto-collected info */}
          <div className="p-2 rounded-md bg-secondary/30 border border-border/50">
            <p className="text-[11px] font-medium mb-1 text-foreground/80">Auto-collected:</p>
            <ul className="text-[10px] text-muted-foreground space-y-0.5 grid grid-cols-2 gap-x-2">
              <li>• Browser & device</li>
              <li>• Page URL</li>
              <li>• Session ID</li>
              <li>• Last 5 messages</li>
            </ul>
          </div>

          {/* Status messages */}
          {submitStatus === "success" && (
            <div className="p-2 rounded-md bg-green-500/10 border border-green-500/20 text-green-600 text-xs flex items-center gap-2">
              <div className="w-1 h-1 rounded-full bg-green-500"></div>
              Thank you! Your bug report has been submitted successfully.
            </div>
          )}

          {submitStatus === "error" && (
            <div className="p-2 rounded-md bg-destructive/10 border border-destructive/20 text-destructive text-xs flex items-center gap-2">
              <div className="w-1 h-1 rounded-full bg-destructive"></div>
              Failed to submit. Please try again or contact bugs@mylittleprice.com
            </div>
          )}
        </form>

        {/* Footer */}
        <div className="px-4 py-3 border-t border-border flex items-center justify-end gap-2 bg-muted/20">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="px-3 py-1.5 rounded-md text-xs font-medium hover:bg-secondary/80 transition-colors disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={isSubmitting || !description.trim()}
            className="px-4 py-1.5 rounded-md text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-1.5"
          >
            {isSubmitting ? (
              <>
                <Loader2 className="w-3.5 h-3.5 animate-spin" />
                <span>Submitting...</span>
              </>
            ) : (
              <>
                <Send className="w-3.5 h-3.5" />
                <span>Submit Report</span>
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );

  // Use portal to render modal at body level
  return createPortal(modalContent, document.body);
}
