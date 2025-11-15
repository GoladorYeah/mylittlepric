"use client";

import { useState, useRef, useEffect, useMemo } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, Globe, Languages, Check, ChevronDown, Moon, Sun, Monitor } from "lucide-react";
import { usePreferences, usePreferenceActions, getCurrencyForCountry, useAuthStore } from "@/shared/lib";
import { useTheme } from "next-themes";
import { COUNTRIES, LANGUAGES } from "@/shared/constants";

export default function SettingsPage() {
  const router = useRouter();
  const { country, language, currency } = usePreferences();
  const { setCountry, setLanguage, setCurrency, syncPreferencesToServer } = usePreferenceActions();
  const { accessToken } = useAuthStore();
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const [countrySearchQuery, setCountrySearchQuery] = useState("");
  const [languageSearchQuery, setLanguageSearchQuery] = useState("");
  const [isCountryDropdownOpen, setIsCountryDropdownOpen] = useState(false);
  const [isLanguageDropdownOpen, setIsLanguageDropdownOpen] = useState(false);

  const countrySearchRef = useRef<HTMLInputElement>(null);
  const languageSearchRef = useRef<HTMLInputElement>(null);

  const selectedCountry = useMemo(
    () => COUNTRIES.find((c) => c.code === country.toLowerCase()) || COUNTRIES[0],
    [country]
  );

  const selectedLanguage = useMemo(
    () => LANGUAGES.find((l) => l.code === language.toLowerCase()) || LANGUAGES[0],
    [language]
  );

  const filteredCountries = useMemo(
    () => COUNTRIES.filter(
      (c) =>
        c.name.toLowerCase().includes(countrySearchQuery.toLowerCase()) ||
        c.code.toLowerCase().includes(countrySearchQuery.toLowerCase())
    ),
    [countrySearchQuery]
  );

  const filteredLanguages = useMemo(
    () => LANGUAGES.filter(
      (l) =>
        l.name.toLowerCase().includes(languageSearchQuery.toLowerCase()) ||
        l.nativeName.toLowerCase().includes(languageSearchQuery.toLowerCase()) ||
        l.code.toLowerCase().includes(languageSearchQuery.toLowerCase())
    ),
    [languageSearchQuery]
  );

  // Set mounted state for theme
  useEffect(() => {
    setMounted(true);
  }, []);

  // Focus search input when dropdown opens
  useEffect(() => {
    if (isCountryDropdownOpen) {
      setTimeout(() => countrySearchRef.current?.focus(), 100);
    }
  }, [isCountryDropdownOpen]);

  useEffect(() => {
    if (isLanguageDropdownOpen) {
      setTimeout(() => languageSearchRef.current?.focus(), 100);
    }
  }, [isLanguageDropdownOpen]);

  const handleCountrySelect = async (countryCode: string) => {
    // Optimistic UI update - update immediately without waiting for server
    setCountry(countryCode);
    const newCurrency = getCurrencyForCountry(countryCode.toUpperCase());
    setCurrency(newCurrency);
    setIsCountryDropdownOpen(false);
    setCountrySearchQuery("");

    // Sync to server in background (non-blocking)
    if (accessToken) {
      syncPreferencesToServer().catch((error) => {
        console.error("Failed to sync country preference:", error);
        // Could show a toast notification here if needed
      });
    }
  };

  const handleLanguageSelect = async (languageCode: string) => {
    // Optimistic UI update - update immediately without waiting for server
    setLanguage(languageCode);
    setIsLanguageDropdownOpen(false);
    setLanguageSearchQuery("");

    // Sync to server in background (non-blocking)
    if (accessToken) {
      syncPreferencesToServer().catch((error) => {
        console.error("Failed to sync language preference:", error);
        // Could show a toast notification here if needed
      });
    }
  };

  const handleThemeChange = async (newTheme: string) => {
    // Optimistic UI update - apply theme immediately
    setTheme(newTheme);

    // Sync theme to server in background (non-blocking)
    if (accessToken) {
      import("@/shared/lib/preferences-api")
        .then(({ PreferencesAPI }) => PreferencesAPI.updateUserPreferences({ theme: newTheme }))
        .then(() => console.log("✅ Synced theme to server:", newTheme))
        .catch((error) => {
          console.error("Failed to sync theme preference:", error);
          // Could show a toast notification here if needed
        });
    }
  };

  const handleBack = () => {
    router.back();
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-10 border-b border-border bg-background">
        <div className="container mx-auto px-4 h-16 flex items-center gap-4">
          <button
            onClick={handleBack}
            className="p-2 hover:bg-secondary rounded-lg transition-colors cursor-pointer"
            aria-label="Go back"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <h1 className="text-2xl font-semibold text-foreground">Settings</h1>
        </div>
      </header>

      {/* Content */}
      <main className="container mx-auto px-4 py-8 max-w-3xl">
        <div className="space-y-8">
          {/* Regional Settings Section */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Globe className="w-5 h-5 text-primary" />
              <h2 className="text-lg font-semibold text-foreground">Regional Settings</h2>
            </div>

            <div className="space-y-4 pl-7">
              {/* Country Selection */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">
                  Country/Region
                </label>
                <div className="relative">
                  <button
                    type="button"
                    onClick={() => setIsCountryDropdownOpen(!isCountryDropdownOpen)}
                    className="w-full flex items-center justify-between px-4 py-3 bg-secondary hover:bg-secondary/80 border border-border rounded-lg transition-colors cursor-pointer"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">{selectedCountry.flag}</span>
                      <span className="text-sm font-medium">{selectedCountry.name}</span>
                    </div>
                    <ChevronDown className={`w-4 h-4 transition-transform ${isCountryDropdownOpen ? 'rotate-180' : ''}`} />
                  </button>

                  {isCountryDropdownOpen && (
                    <div className="absolute top-full left-0 right-0 mt-2 bg-background border border-border rounded-lg shadow-xl overflow-hidden z-50">
                      {/* Search Input */}
                      <div className="p-3 border-b border-border">
                        <input
                          ref={countrySearchRef}
                          type="text"
                          value={countrySearchQuery}
                          onChange={(e) => setCountrySearchQuery(e.target.value)}
                          placeholder="Search countries..."
                          className="w-full px-3 py-2 rounded-md bg-secondary border border-border focus:border-primary focus:outline-none transition-colors text-sm"
                        />
                      </div>

                      {/* Countries List */}
                      <div className="max-h-64 overflow-y-auto">
                        {filteredCountries.length > 0 ? (
                          filteredCountries.map((c) => (
                            <button
                              key={c.code}
                              onClick={() => handleCountrySelect(c.code)}
                              className={`w-full px-4 py-2.5 flex items-center justify-between hover:bg-secondary transition-colors text-left cursor-pointer ${
                                c.code === country.toLowerCase() ? "bg-secondary/50" : ""
                              }`}
                            >
                              <div className="flex items-center gap-3">
                                <span className="text-xl">{c.flag}</span>
                                <span className="text-sm font-medium">{c.name}</span>
                              </div>
                              {c.code === country.toLowerCase() && (
                                <Check className="w-4 h-4 text-primary" />
                              )}
                            </button>
                          ))
                        ) : (
                          <div className="px-4 py-8 text-center text-sm text-muted-foreground">
                            No countries found
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
                <p className="text-xs text-muted-foreground">
                  Currency: {currency}
                </p>
              </div>
            </div>
          </div>

          {/* Language Settings Section */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Languages className="w-5 h-5 text-primary" />
              <h2 className="text-lg font-semibold text-foreground">Language Settings</h2>
            </div>

            <div className="space-y-4 pl-7">
              {/* Language Selection */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">
                  Agent Communication Language
                </label>
                <p className="text-xs text-muted-foreground">
                  Choose the language for the AI assistant to communicate with you
                </p>
                <div className="relative">
                  <button
                    type="button"
                    onClick={() => setIsLanguageDropdownOpen(!isLanguageDropdownOpen)}
                    className="w-full flex items-center justify-between px-4 py-3 bg-secondary hover:bg-secondary/80 border border-border rounded-lg transition-colors cursor-pointer"
                  >
                    <div className="flex items-center gap-3">
                      <Languages className="w-5 h-5 text-muted-foreground" />
                      <div className="flex flex-col items-start">
                        <span className="text-sm font-medium">{selectedLanguage.name}</span>
                        <span className="text-xs text-muted-foreground">{selectedLanguage.nativeName}</span>
                      </div>
                    </div>
                    <ChevronDown className={`w-4 h-4 transition-transform ${isLanguageDropdownOpen ? 'rotate-180' : ''}`} />
                  </button>

                  {isLanguageDropdownOpen && (
                    <div className="absolute top-full left-0 right-0 mt-2 bg-background border border-border rounded-lg shadow-xl overflow-hidden z-50">
                      {/* Search Input */}
                      <div className="p-3 border-b border-border">
                        <input
                          ref={languageSearchRef}
                          type="text"
                          value={languageSearchQuery}
                          onChange={(e) => setLanguageSearchQuery(e.target.value)}
                          placeholder="Search languages..."
                          className="w-full px-3 py-2 rounded-md bg-secondary border border-border focus:border-primary focus:outline-none transition-colors text-sm"
                        />
                      </div>

                      {/* Languages List */}
                      <div className="max-h-64 overflow-y-auto">
                        {filteredLanguages.length > 0 ? (
                          filteredLanguages.map((l) => (
                            <button
                              key={l.code}
                              onClick={() => handleLanguageSelect(l.code)}
                              className={`w-full px-4 py-2.5 flex items-center justify-between hover:bg-secondary transition-colors text-left cursor-pointer ${
                                l.code === language.toLowerCase() ? "bg-secondary/50" : ""
                              }`}
                            >
                              <div className="flex flex-col">
                                <span className="text-sm font-medium">{l.name}</span>
                                <span className="text-xs text-muted-foreground">{l.nativeName}</span>
                              </div>
                              {l.code === language.toLowerCase() && (
                                <Check className="w-4 h-4 text-primary" />
                              )}
                            </button>
                          ))
                        ) : (
                          <div className="px-4 py-8 text-center text-sm text-muted-foreground">
                            No languages found
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>

          {/* Appearance Settings Section */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Monitor className="w-5 h-5 text-primary" />
              <h2 className="text-lg font-semibold text-foreground">Appearance</h2>
            </div>

            <div className="space-y-4 pl-7">
              {/* Theme Selection */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-muted-foreground">
                  Theme
                </label>
                <p className="text-xs text-muted-foreground">
                  Choose your preferred color scheme
                </p>
                {mounted && (
                  <div className="grid grid-cols-3 gap-3 pt-2">
                    <button
                      onClick={() => handleThemeChange("light")}
                      className={`flex flex-col items-center gap-2 p-4 rounded-lg border transition-all cursor-pointer ${
                        theme === "light"
                          ? "bg-primary/10 border-primary"
                          : "bg-secondary border-border hover:border-primary/50"
                      }`}
                    >
                      <Sun className={`w-6 h-6 ${theme === "light" ? "text-primary" : "text-muted-foreground"}`} />
                      <span className="text-sm font-medium">Light</span>
                      {theme === "light" && <Check className="w-4 h-4 text-primary" />}
                    </button>

                    <button
                      onClick={() => handleThemeChange("dark")}
                      className={`flex flex-col items-center gap-2 p-4 rounded-lg border transition-all cursor-pointer ${
                        theme === "dark"
                          ? "bg-primary/10 border-primary"
                          : "bg-secondary border-border hover:border-primary/50"
                      }`}
                    >
                      <Moon className={`w-6 h-6 ${theme === "dark" ? "text-primary" : "text-muted-foreground"}`} />
                      <span className="text-sm font-medium">Dark</span>
                      {theme === "dark" && <Check className="w-4 h-4 text-primary" />}
                    </button>

                    <button
                      onClick={() => handleThemeChange("system")}
                      className={`flex flex-col items-center gap-2 p-4 rounded-lg border transition-all cursor-pointer ${
                        theme === "system"
                          ? "bg-primary/10 border-primary"
                          : "bg-secondary border-border hover:border-primary/50"
                      }`}
                    >
                      <Monitor className={`w-6 h-6 ${theme === "system" ? "text-primary" : "text-muted-foreground"}`} />
                      <span className="text-sm font-medium">System</span>
                      {theme === "system" && <Check className="w-4 h-4 text-primary" />}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
