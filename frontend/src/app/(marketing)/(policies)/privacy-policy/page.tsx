import { PolicyLayout } from "@/components/PolicyLayout";
import { PrivacyPolicyContent } from "@/components/policies/PrivacyPolicyContent";

export const metadata = {
  title: "Privacy Policy - MyLittlePrice",
  description: "MyLittlePrice Privacy Policy - Learn how we collect, use, and protect your personal information",
};

export default function PrivacyPolicyPage() {
  return (
    <PolicyLayout title="Privacy Policy" lastUpdated="November 7, 2025">
      <PrivacyPolicyContent />
    </PolicyLayout>
  );
}
