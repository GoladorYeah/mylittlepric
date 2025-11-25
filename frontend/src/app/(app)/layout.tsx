"use client";

// Note: This layout no longer requires authentication
// Anonymous users can now access the chat with limited searches (3 free searches)
// After reaching the limit, they will be prompted to sign up/log in

export default function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return children;
}
