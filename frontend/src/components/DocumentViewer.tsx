"use client";

import { useRouter } from "next/navigation";
import { ArrowLeft, Download } from "lucide-react";
import { Logo } from "@/components/Logo";

interface DocumentViewerProps {
  title: string;
  pdfPath: string;
  lastUpdated: string;
}

export function DocumentViewer({ title, pdfPath, lastUpdated }: DocumentViewerProps) {
  const router = useRouter();

  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
      <header className="fixed top-0 left-0 right-0 z-50 bg-background border-b border-border">
        <div className="container mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <button
              onClick={() => router.back()}
              className="flex items-center gap-2 text-muted-foreground hover:text-foreground transition-colors"
            >
              <ArrowLeft className="w-5 h-5" />
              <span className="hidden sm:inline">Back</span>
            </button>
            <div className="h-6 w-px bg-border hidden sm:block" />
            <Logo width={84.24} height={32} />
          </div>

          <a
            href={pdfPath}
            download
            className="flex items-center gap-2 px-4 py-2 bg-primary text-primary-foreground rounded-full font-semibold hover:opacity-90 transition-opacity text-sm"
          >
            <Download className="w-4 h-4" />
            <span className="hidden sm:inline">Download PDF</span>
            <span className="sm:hidden">Download</span>
          </a>
        </div>
      </header>

      <main className="pt-24 pb-12">
        <div className="container mx-auto px-4">
          <div className="max-w-5xl mx-auto space-y-6">
            <div className="text-center space-y-2">
              <h1 className="text-3xl md:text-4xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                {title}
              </h1>
              <p className="text-sm text-muted-foreground">
                Last updated: {lastUpdated}
              </p>
            </div>

            <div className="bg-secondary/50 border border-border rounded-2xl overflow-hidden shadow-xl">
              <div className="w-full h-[calc(100vh-220px)] min-h-[600px]">
                <iframe
                  src={pdfPath}
                  className="w-full h-full"
                  title={title}
                />
              </div>
            </div>

            <div className="text-center">
              <p className="text-sm text-muted-foreground">
                If the document doesn't load,{" "}
                <a
                  href={pdfPath}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline"
                >
                  click here to open it in a new tab
                </a>
              </p>
            </div>
          </div>
        </div>
      </main>

      <footer className="border-t border-border py-8 mt-12">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-muted-foreground">
            <p>&copy; 2025 MyLittlePrice. All rights reserved.</p>
            <div className="flex gap-6">
              <a href="/privacy-policy" className="hover:text-foreground transition-colors">
                Privacy Policy
              </a>
              <a href="/terms-of-use" className="hover:text-foreground transition-colors">
                Terms of Use
              </a>
              <a href="/cookie-policy" className="hover:text-foreground transition-colors">
                Cookie Policy
              </a>
              <a href="/advertising-policy" className="hover:text-foreground transition-colors">
                Advertising Policy
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
