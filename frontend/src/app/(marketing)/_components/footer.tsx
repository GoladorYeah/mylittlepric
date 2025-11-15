import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t border-border py-12">
      <div className="container mx-auto px-4">
        <div className="grid md:grid-cols-4 gap-8 mb-8">
          {/* About Section */}
          <div className="space-y-3">
            <h3 className="font-semibold text-foreground">MyLittlePrice</h3>
            <p className="text-sm text-muted-foreground">
              AI-powered shopping assistant helping you find the best products at the best prices.
            </p>
          </div>

          {/* Quick Links */}
          <div className="space-y-3">
            <h3 className="font-semibold text-foreground">Quick Links</h3>
            <nav className="flex flex-col gap-2">
              <Link
                href="/about"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                About Us
              </Link>
              <Link
                href="/contact"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Contact
              </Link>
            </nav>
          </div>

          {/* Legal */}
          <div className="space-y-3">
            <h3 className="font-semibold text-foreground">Legal</h3>
            <nav className="flex flex-col gap-2">
              <Link
                href="/privacy-policy"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Privacy Policy
              </Link>
              <Link
                href="/terms-of-use"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Terms of Use
              </Link>
              <Link
                href="/disclaimer"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Disclaimer
              </Link>
            </nav>
          </div>

          {/* Policies */}
          <div className="space-y-3">
            <h3 className="font-semibold text-foreground">Policies</h3>
            <nav className="flex flex-col gap-2">
              <Link
                href="/cookie-policy"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Cookie Policy
              </Link>
              <Link
                href="/advertising-policy"
                className="text-sm text-muted-foreground hover:text-foreground transition-colors"
              >
                Advertising Policy
              </Link>
            </nav>
          </div>
        </div>

        {/* Copyright */}
        <div className="pt-8 border-t border-border text-center text-sm text-muted-foreground">
          <p>&copy; {new Date().getFullYear()} MyLittlePrice. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
