import { ProductDetailsResponse, ChatMessage } from "@/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function getProductDetails(
  pageToken: string,
  country: string
): Promise<ProductDetailsResponse> {
  const response = await fetch(`${API_URL}/api/product-details`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ page_token: pageToken, country }),
  });

  if (!response.ok) {
    throw new Error("Failed to fetch product details");
  }

  return response.json();
}

export interface SessionMessagesResponse {
  messages: Array<{
    role: string;
    content: string;
    timestamp?: string;
    quick_replies?: string[];
    products?: any[];
    search_type?: string;
  }>;
  session_id: string;
  message_count: number;
  search_state?: {
    status: string;
    category: string;
    search_count: number;
    last_product?: any;
  };
}

export async function getSessionMessages(
  sessionId: string
): Promise<SessionMessagesResponse> {
  const response = await fetch(
    `${API_URL}/api/chat/messages?session_id=${encodeURIComponent(sessionId)}`,
    {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (!response.ok) {
    throw new Error("Failed to fetch session messages");
  }

  return response.json();
}