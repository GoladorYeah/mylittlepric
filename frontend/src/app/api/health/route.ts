import { NextResponse } from "next/server";

export async function GET() {
  try {
    const backendUrl =
      process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const response = await fetch(`${backendUrl}/health`, {
      cache: "no-store",
    });

    if (!response.ok) {
      return NextResponse.json(
        { status: "unhealthy", backend: "down" },
        { status: 503 }
      );
    }

    const data = await response.json();
    return NextResponse.json(
      {
        status: "healthy",
        backend: data.status || "up",
        timestamp: new Date().toISOString(),
      },
      { status: 200 }
    );
  } catch (error) {
    return NextResponse.json(
      {
        status: "unhealthy",
        backend: "unreachable",
        error: error instanceof Error ? error.message : "Unknown error",
      },
      { status: 503 }
    );
  }
}

