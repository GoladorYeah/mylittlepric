import { DocumentViewer } from "@/components/DocumentViewer";

export const metadata = {
  title: "Terms of Use - MyLittlePrice",
  description: "MyLittlePrice Terms of Use - Read our terms and conditions for using our service",
};

export default function TermsOfUsePage() {
  return (
    <DocumentViewer
      title="Terms of Use"
      pdfPath="/documents/terms-of-use.pdf"
      lastUpdated="November 6, 2025"
    />
  );
}
