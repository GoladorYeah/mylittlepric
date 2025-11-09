import Link from "next/link";

export function Footer() {
  return (
    <footer className="border-t border-border py-8">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-muted-foreground">
          <p>&copy; 2025 MyLittlePrice. All rights reserved.</p>
          <nav className="flex gap-6">
            <Link
              href="/privacy-policy"
              className="hover:text-foreground transition-colors"
            >
              Privacy Policy
            </Link>
            <Link
              href="/terms-of-use"
              className="hover:text-foreground transition-colors"
            >
              Terms of Use
            </Link>
            <Link
              href="/cookie-policy"
              className="hover:text-foreground transition-colors"
            >
              Cookie Policy
            </Link>
            <Link
              href="/advertising-policy"
              className="hover:text-foreground transition-colors"
            >
              Advertising Policy
            </Link>
          </nav>
        </div>
      </div>
    </footer>
  );
}
