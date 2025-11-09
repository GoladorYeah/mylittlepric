import { DocumentViewer } from "@/components/DocumentViewer";

export const metadata = {
  title: "Advertising Policy - MyLittlePrice",
  description: "MyLittlePrice Advertiser & Seller Advertising Policy - Guidelines for advertisers and sellers",
};

export default function AdvertisingPolicyPage() {
  return (
    <DocumentViewer
      title="Advertiser & Seller Advertising Policy"
      pdfPath="/documents/advertising-policy.pdf"
      lastUpdated="November 6, 2025"
    />
  );
}
