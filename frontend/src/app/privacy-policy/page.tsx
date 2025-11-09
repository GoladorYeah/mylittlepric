import { DocumentViewer } from "@/components/DocumentViewer";

export const metadata = {
  title: "Privacy Policy - MyLittlePrice",
  description: "MyLittlePrice Privacy Policy - Learn how we collect, use, and protect your personal information",
};

export default function PrivacyPolicyPage() {
  return (
    <DocumentViewer
      title="Privacy Policy"
      pdfPath="/documents/privacy-policy.pdf"
      lastUpdated="November 7, 2025"
    />
  );
}
