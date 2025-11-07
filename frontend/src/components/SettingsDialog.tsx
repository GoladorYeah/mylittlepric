"use client";

import { useState, useRef, useEffect } from "react";
import { X, Globe, Languages, Check, ChevronDown } from "lucide-react";
import { useChatStore } from "@/lib/store";
import { getCurrencyForCountry } from "@/lib/locale";

interface Country {
  code: string;
  name: string;
  flag: string;
}

interface Language {
  code: string;
  name: string;
  nativeName: string;
}

const COUNTRIES: Country[] = [
  { code: "us", name: "United States", flag: "🇺🇸" },
  { code: "gb", name: "United Kingdom", flag: "🇬🇧" },
  { code: "ca", name: "Canada", flag: "🇨🇦" },
  { code: "au", name: "Australia", flag: "🇦🇺" },
  { code: "de", name: "Germany", flag: "🇩🇪" },
  { code: "fr", name: "France", flag: "🇫🇷" },
  { code: "es", name: "Spain", flag: "🇪🇸" },
  { code: "it", name: "Italy", flag: "🇮🇹" },
  { code: "nl", name: "Netherlands", flag: "🇳🇱" },
  { code: "be", name: "Belgium", flag: "🇧🇪" },
  { code: "ch", name: "Switzerland", flag: "🇨🇭" },
  { code: "at", name: "Austria", flag: "🇦🇹" },
  { code: "se", name: "Sweden", flag: "🇸🇪" },
  { code: "no", name: "Norway", flag: "🇳🇴" },
  { code: "dk", name: "Denmark", flag: "🇩🇰" },
  { code: "fi", name: "Finland", flag: "🇫🇮" },
  { code: "pl", name: "Poland", flag: "🇵🇱" },
  { code: "cz", name: "Czech Republic", flag: "🇨🇿" },
  { code: "pt", name: "Portugal", flag: "🇵🇹" },
  { code: "gr", name: "Greece", flag: "🇬🇷" },
  { code: "ie", name: "Ireland", flag: "🇮🇪" },
  { code: "jp", name: "Japan", flag: "🇯🇵" },
  { code: "kr", name: "South Korea", flag: "🇰🇷" },
  { code: "cn", name: "China", flag: "🇨🇳" },
  { code: "in", name: "India", flag: "🇮🇳" },
  { code: "sg", name: "Singapore", flag: "🇸🇬" },
  { code: "hk", name: "Hong Kong", flag: "🇭🇰" },
  { code: "tw", name: "Taiwan", flag: "🇹🇼" },
  { code: "nz", name: "New Zealand", flag: "🇳🇿" },
  { code: "mx", name: "Mexico", flag: "🇲🇽" },
  { code: "br", name: "Brazil", flag: "🇧🇷" },
  { code: "ar", name: "Argentina", flag: "🇦🇷" },
  { code: "cl", name: "Chile", flag: "🇨🇱" },
  { code: "za", name: "South Africa", flag: "🇿🇦" },
  { code: "ae", name: "UAE", flag: "🇦🇪" },
  { code: "sa", name: "Saudi Arabia", flag: "🇸🇦" },
  { code: "tr", name: "Turkey", flag: "🇹🇷" },
  { code: "ru", name: "Russia", flag: "🇷🇺" },
  { code: "ua", name: "Ukraine", flag: "🇺🇦" },
  { code: "il", name: "Israel", flag: "🇮🇱" },
  { code: "eg", name: "Egypt", flag: "🇪🇬" },
  { code: "th", name: "Thailand", flag: "🇹🇭" },
  { code: "my", name: "Malaysia", flag: "🇲🇾" },
  { code: "id", name: "Indonesia", flag: "🇮🇩" },
  { code: "ph", name: "Philippines", flag: "🇵🇭" },
  { code: "vn", name: "Vietnam", flag: "🇻🇳" },
];

const LANGUAGES: Language[] = [
  { code: "en", name: "English", nativeName: "English" },
  { code: "es", name: "Spanish", nativeName: "Español" },
  { code: "fr", name: "French", nativeName: "Français" },
  { code: "de", name: "German", nativeName: "Deutsch" },
  { code: "it", name: "Italian", nativeName: "Italiano" },
  { code: "pt", name: "Portuguese", nativeName: "Português" },
  { code: "ru", name: "Russian", nativeName: "Русский" },
  { code: "zh", name: "Chinese", nativeName: "中文" },
  { code: "ja", name: "Japanese", nativeName: "日本語" },
  { code: "ko", name: "Korean", nativeName: "한국어" },
  { code: "ar", name: "Arabic", nativeName: "العربية" },
  { code: "hi", name: "Hindi", nativeName: "हिन्दी" },
  { code: "nl", name: "Dutch", nativeName: "Nederlands" },
  { code: "pl", name: "Polish", nativeName: "Polski" },
  { code: "tr", name: "Turkish", nativeName: "Türkçe" },
  { code: "sv", name: "Swedish", nativeName: "Svenska" },
  { code: "no", name: "Norwegian", nativeName: "Norsk" },
  { code: "da", name: "Danish", nativeName: "Dansk" },
  { code: "fi", name: "Finnish", nativeName: "Suomi" },
  { code: "cs", name: "Czech", nativeName: "Čeština" },
  { code: "el", name: "Greek", nativeName: "Ελληνικά" },
  { code: "he", name: "Hebrew", nativeName: "עברית" },
  { code: "th", name: "Thai", nativeName: "ไทย" },
  { code: "vi", name: "Vietnamese", nativeName: "Tiếng Việt" },
  { code: "id", name: "Indonesian", nativeName: "Bahasa Indonesia" },
  { code: "ms", name: "Malay", nativeName: "Bahasa Melayu" },
  { code: "uk", name: "Ukrainian", nativeName: "Українська" },
  { code: "ro", name: "Romanian", nativeName: "Română" },
  { code: "hu", name: "Hungarian", nativeName: "Magyar" },
  { code: "hr", name: "Croatian", nativeName: "Hrvatski" },
];

interface SettingsDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function SettingsDialog({ isOpen, onClose }: SettingsDialogProps) {
  const { country, language, currency, setCountry, setLanguage, setCurrency } = useChatStore();
  const [countrySearchQuery, setCountrySearchQuery] = useState("");
  const [languageSearchQuery, setLanguageSearchQuery] = useState("");
  const [isCountryDropdownOpen, setIsCountryDropdownOpen] = useState(false);
  const [isLanguageDropdownOpen, setIsLanguageDropdownOpen] = useState(false);

  const dialogRef = useRef<HTMLDivElement>(null);
  const countrySearchRef = useRef<HTMLInputElement>(null);
  const languageSearchRef = useRef<HTMLInputElement>(null);

  const selectedCountry = COUNTRIES.find((c) => c.code === country.toLowerCase()) || COUNTRIES[0];
  const selectedLanguage = LANGUAGES.find((l) => l.code === language.toLowerCase()) || LANGUAGES[0];

  const filteredCountries = COUNTRIES.filter(
    (c) =>
      c.name.toLowerCase().includes(countrySearchQuery.toLowerCase()) ||
      c.code.toLowerCase().includes(countrySearchQuery.toLowerCase())
  );

  const filteredLanguages = LANGUAGES.filter(
    (l) =>
      l.name.toLowerCase().includes(languageSearchQuery.toLowerCase()) ||
      l.nativeName.toLowerCase().includes(languageSearchQuery.toLowerCase()) ||
      l.code.toLowerCase().includes(languageSearchQuery.toLowerCase())
  );

  // Close dialog on ESC key
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape" && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    };
  }, [isOpen, onClose]);

  // Close dialog when clicking outside
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (dialogRef.current && !dialogRef.current.contains(e.target as Node)) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen, onClose]);

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

  const handleCountrySelect = (countryCode: string) => {
    setCountry(countryCode);
    const newCurrency = getCurrencyForCountry(countryCode.toUpperCase());
    setCurrency(newCurrency);
    setIsCountryDropdownOpen(false);
    setCountrySearchQuery("");
  };

  const handleLanguageSelect = (languageCode: string) => {
    setLanguage(languageCode);
    setIsLanguageDropdownOpen(false);
    setLanguageSearchQuery("");
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div
        ref={dialogRef}
        className="relative w-full max-w-2xl max-h-[90vh] overflow-y-auto bg-background rounded-lg shadow-2xl border border-border m-4"
      >
        {/* Header */}
        <div className="sticky top-0 z-10 flex items-center justify-between p-6 border-b border-border bg-background">
          <h2 className="text-2xl font-semibold text-foreground">Settings</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-secondary rounded-lg transition-colors cursor-pointer"
            aria-label="Close settings"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-8">
          {/* Regional Settings Section */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Globe className="w-5 h-5 text-primary" />
              <h3 className="text-lg font-semibold text-foreground">Regional Settings</h3>
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
              <h3 className="text-lg font-semibold text-foreground">Language Settings</h3>
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
        </div>

        {/* Footer */}
        <div className="sticky bottom-0 flex items-center justify-end gap-3 p-6 border-t border-border bg-background">
          <button
            onClick={onClose}
            className="px-6 py-2 bg-primary text-primary-foreground rounded-lg hover:opacity-90 transition-opacity cursor-pointer font-medium"
          >
            Done
          </button>
        </div>
      </div>
    </div>
  );
}
