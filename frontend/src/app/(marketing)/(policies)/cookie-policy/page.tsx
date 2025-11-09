import { PolicyLayout, CookiePolicyContent } from "@/features/policies";

export const metadata = {
  title: "Cookie Policy - MyLittlePrice",
  description: "MyLittlePrice Cookie Policy - Learn about how we use cookies and similar technologies",
};

export default function CookiePolicyPage() {
  return (
    <PolicyLayout title="Cookie Policy" lastUpdated="November 7, 2025">
      <CookiePolicyContent />
    </PolicyLayout>
  );
}
