"use client";

import Image from "next/image";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

interface LogoProps {
  className?: string;
  width?: number;
  height?: number;
}

export function Logo({ className = "", width = 84.24, height = 32 }: LogoProps) {
  const { resolvedTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  // Prevent hydration mismatch
  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    // Return a placeholder during SSR to avoid hydration mismatch
    return (
      <div
        className={className}
        style={{ width: `${width}px`, height: `${height}px` }}
      />
    );
  }

  const logoSrc =
    resolvedTheme === "dark" ? "/dark-logo.svg" : "/light-logo.svg";

  return (
    <Image
      src={logoSrc}
      alt="MyLittlePrice Logo"
      width={width}
      height={height}
      className={className}
      priority
    />
  );
}
