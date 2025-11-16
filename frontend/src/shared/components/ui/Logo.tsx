"use client";

import Image from "next/image";

interface LogoProps {
  className?: string;
  width?: number;
  height?: number;
}

export function Logo({ className = "", width = 120, height = 40 }: LogoProps) {
  return (
    <Image
      src="/logo.svg"
      alt="MyLittlePrice Logo"
      width={width}
      height={height}
      className={className}
      priority
    />
  );
}
