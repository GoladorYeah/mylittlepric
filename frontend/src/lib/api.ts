import { ProductDetailsResponse } from "@/types";

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