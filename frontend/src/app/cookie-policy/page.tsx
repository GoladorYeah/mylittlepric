import { DocumentViewer } from "@/components/DocumentViewer";

export const metadata = {
  title: "Cookie Policy - MyLittlePrice",
  description: "MyLittlePrice Cookie Policy - Learn about how we use cookies and similar technologies",
};

export default function CookiePolicyPage() {
  return (
    <DocumentViewer
      title="Cookie Policy"
      pdfPath="/documents/cookie-policy.pdf"
      lastUpdated="November 7, 2025"
    />
  );
}
