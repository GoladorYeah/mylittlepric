import type { Metadata } from "next";
import { PolicyLayout } from "@/features/policies";

export const metadata: Metadata = {
  title: "Disclaimer - MyLittlePrice",
  description: "Important disclaimers and limitations regarding the use of MyLittlePrice services.",
};

export default function DisclaimerPage() {
  return (
    <PolicyLayout
      title="Disclaimer"
      lastUpdated={new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
    >
      <section className="space-y-4">
        <h2 className="text-3xl font-bold">General Information</h2>
        <p className="text-muted-foreground leading-relaxed">
          The information provided by MyLittlePrice ("we," "us," or "our")
          on our website and through our services is for general
          informational purposes only. All information on the site is
          provided in good faith, however we make no representation or
          warranty of any kind, express or implied, regarding the accuracy,
          adequacy, validity, reliability, availability, or completeness of
          any information on the site.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">AI-Generated Recommendations</h2>
        <p className="text-muted-foreground leading-relaxed">
          MyLittlePrice uses artificial intelligence (AI) technology to
          provide product recommendations and shopping assistance. Please
          be aware that:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>
            AI-generated recommendations are based on algorithms and may
            not always reflect the most suitable choice for your specific
            needs
          </li>
          <li>
            The AI may occasionally provide inaccurate or incomplete
            information about products
          </li>
          <li>
            Product descriptions, prices, and availability are sourced from
            third-party retailers and may not always be current or accurate
          </li>
          <li>
            We recommend verifying all product information directly with
            the retailer before making a purchase
          </li>
        </ul>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Product Information and Pricing</h2>
        <p className="text-muted-foreground leading-relaxed">
          All product information, including prices, descriptions,
          specifications, and availability displayed on MyLittlePrice is:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>
            Sourced from third-party retailers and search APIs
          </li>
          <li>
            Subject to change without notice
          </li>
          <li>
            Not guaranteed to be accurate, complete, or current at the time
            of viewing
          </li>
          <li>
            Dependent on the accuracy of information provided by retailers
            and manufacturers
          </li>
        </ul>
        <p className="text-muted-foreground leading-relaxed">
          We are not responsible for any discrepancies between the
          information displayed on our platform and the actual product
          information on retailer websites.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Affiliate Relationships</h2>
        <p className="text-muted-foreground leading-relaxed">
          MyLittlePrice may contain affiliate links to third-party products
          and services. This means:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>
            We may earn a commission when you click on certain links or
            make purchases through our platform
          </li>
          <li>
            These affiliate relationships do not influence our product
            recommendations or search results
          </li>
          <li>
            We strive to maintain objectivity in all product suggestions
          </li>
          <li>
            Any commissions earned help us maintain and improve our service
          </li>
        </ul>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Third-Party Websites</h2>
        <p className="text-muted-foreground leading-relaxed">
          Our service contains links to third-party websites and retailers.
          We have no control over, and assume no responsibility for:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>
            The content, privacy policies, or practices of any third-party
            sites or services
          </li>
          <li>
            The quality, safety, or legality of products sold by third-party
            retailers
          </li>
          <li>
            Any transactions conducted between you and third-party retailers
          </li>
          <li>
            Shipping, returns, refunds, or customer service provided by
            third parties
          </li>
        </ul>
        <p className="text-muted-foreground leading-relaxed">
          You acknowledge and agree that we shall not be responsible or
          liable, directly or indirectly, for any damage or loss caused or
          alleged to be caused by or in connection with the use of any such
          content, goods, or services available on or through any such
          websites or services.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">No Professional Advice</h2>
        <p className="text-muted-foreground leading-relaxed">
          The information provided through MyLittlePrice is not intended to
          be and does not constitute:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>Financial, investment, or professional purchasing advice</li>
          <li>Medical, health, or safety recommendations</li>
          <li>Legal advice or guidance</li>
          <li>Expert opinion on product quality or suitability</li>
        </ul>
        <p className="text-muted-foreground leading-relaxed">
          You should consult with appropriate professionals before making
          any purchasing decisions based on information provided through our
          service.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Limitation of Liability</h2>
        <p className="text-muted-foreground leading-relaxed">
          Under no circumstance shall we have any liability to you for any
          loss or damage of any kind incurred as a result of:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>
            The use of the site or reliance on any information provided on
            the site
          </li>
          <li>
            Purchases made based on AI recommendations or product
            information displayed
          </li>
          <li>Inaccurate pricing or product information</li>
          <li>
            Transactions with third-party retailers accessed through our
            platform
          </li>
          <li>
            Any interruption, suspension, or termination of our service
          </li>
        </ul>
        <p className="text-muted-foreground leading-relaxed">
          Your use of the site and your reliance on any information on the
          site is solely at your own risk.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Product Returns and Refunds</h2>
        <p className="text-muted-foreground leading-relaxed">
          MyLittlePrice does not sell products directly and is not
          responsible for:
        </p>
        <ul className="list-disc list-inside space-y-2 text-muted-foreground ml-4">
          <li>Processing returns or refunds for products purchased</li>
          <li>Mediating disputes between customers and retailers</li>
          <li>Warranty claims or product defects</li>
          <li>Shipping delays or damages</li>
        </ul>
        <p className="text-muted-foreground leading-relaxed">
          All returns, refunds, and customer service issues must be handled
          directly with the retailer from whom you purchased the product.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Changes to This Disclaimer</h2>
        <p className="text-muted-foreground leading-relaxed">
          We may update this Disclaimer from time to time. We will notify
          you of any changes by posting the new Disclaimer on this page and
          updating the "Last updated" date.
        </p>
        <p className="text-muted-foreground leading-relaxed">
          You are advised to review this Disclaimer periodically for any
          changes. Changes to this Disclaimer are effective when they are
          posted on this page.
        </p>
      </section>

      <section className="space-y-4">
        <h2 className="text-3xl font-bold">Contact Information</h2>
        <p className="text-muted-foreground leading-relaxed">
          If you have any questions about this Disclaimer, please contact
          us:
        </p>
        <div className="p-6 rounded-lg border border-border bg-secondary/20">
          <p className="text-muted-foreground">
            <strong>Email:</strong>{" "}
            <a
              href="mailto:legal@mylittleprice.com"
              className="text-primary hover:underline"
            >
              legal@mylittleprice.com
            </a>
          </p>
          <p className="text-muted-foreground mt-2">
            <strong>Contact Page:</strong>{" "}
            <a href="/contact" className="text-primary hover:underline">
              mylittleprice.com/contact
            </a>
          </p>
        </div>
      </section>

      <section className="p-6 rounded-lg border border-primary/20 bg-primary/5">
        <p className="text-sm text-muted-foreground">
          <strong className="text-foreground">Important Notice:</strong> By
          using MyLittlePrice, you acknowledge that you have read,
          understood, and agree to be bound by this Disclaimer. If you do
          not agree with any part of this Disclaimer, please discontinue
          use of our service immediately.
        </p>
      </section>
    </PolicyLayout>
  );
}
