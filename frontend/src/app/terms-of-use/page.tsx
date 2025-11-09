import { PolicyLayout } from "@/components/PolicyLayout";
import { TermsOfUseContent } from "@/components/policies/TermsOfUseContent";

export const metadata = {
  title: "Terms of Use - MyLittlePrice",
  description: "MyLittlePrice Terms of Use - Read our terms and conditions for using our service",
};

export default function TermsOfUsePage() {
  return (
    <PolicyLayout title="Terms of Use" lastUpdated="November 6, 2025">
      <TermsOfUseContent />
    </PolicyLayout>
  );
}
