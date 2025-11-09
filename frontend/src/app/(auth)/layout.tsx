import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Login - MyLittlePrice",
  description: "Sign in to your MyLittlePrice account",
};

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      {children}
    </div>
  );
}
