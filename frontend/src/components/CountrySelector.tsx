"use client";

import { useState, useRef, useEffect } from "react";
import { Globe, Check } from "lucide-react";
import { useChatStore } from "@/lib/store";

interface Country {
  code: string;
  name: string;
  flag: string;
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

export function CountrySelector() {
  const { country, setCountry } = useChatStore();
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const dropdownRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const selectedCountry = COUNTRIES.find((c) => c.code === country.toLowerCase()) || COUNTRIES[0];

  const filteredCountries = COUNTRIES.filter(
    (c) =>
      c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.code.toLowerCase().includes(searchQuery.toLowerCase())
  );

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setSearchQuery("");
      }
    };

    if (isOpen) {
      document.addEventListener("mousedown", handleClickOutside);
      // Focus search input when dropdown opens
      setTimeout(() => searchInputRef.current?.focus(), 100);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isOpen]);

  const handleCountrySelect = (countryCode: string) => {
    setCountry(countryCode);
    setIsOpen(false);
    setSearchQuery("");
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-1.5 px-3 py-2 rounded-lg bg-secondary hover:bg-secondary/80 transition-colors border border-border"
        title="Select country"
      >
        <Globe className="w-4 h-4 text-muted-foreground" />
        <span className="text-lg">{selectedCountry.flag}</span>
      </button>

      {isOpen && (
        <div className="absolute left-0 bottom-full mb-2 w-72 bg-background border border-border rounded-lg shadow-lg overflow-hidden z-50">
          {/* Search Input */}
          <div className="p-3 border-b border-border">
            <input
              ref={searchInputRef}
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
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
                  className={`w-full px-4 py-2.5 flex items-center justify-between hover:bg-secondary transition-colors text-left ${
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
  );
}
