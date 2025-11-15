/**
 * Shared country data for country selectors across the application
 * This prevents duplication and improves maintainability
 */

export interface Country {
  code: string;
  name: string;
  flag: string;
}

export const COUNTRIES: Country[] = [
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

/**
 * Find a country by its code (case-insensitive)
 */
export function findCountryByCode(code: string): Country | undefined {
  return COUNTRIES.find((c) => c.code === code.toLowerCase());
}

/**
 * Get the default country (US)
 */
export function getDefaultCountry(): Country {
  return COUNTRIES[0];
}
