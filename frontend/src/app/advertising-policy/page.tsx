import { PolicyLayout } from "@/components/PolicyLayout";
import { AdvertisingPolicyContent } from "@/components/policies/AdvertisingPolicyContent";

export const metadata = {
  title: "Advertising Policy - MyLittlePrice",
  description: "MyLittlePrice Advertiser & Seller Advertising Policy - Guidelines for advertisers and sellers",
};

export default function AdvertisingPolicyPage() {
  return (
    <PolicyLayout title="Advertiser & Seller Advertising Policy" lastUpdated="November 6, 2025">
      <AdvertisingPolicyContent />
    </PolicyLayout>
  );
}
