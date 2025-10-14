export async function detectCountryByBackend(): Promise<string | null> {
  try {
    const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const response = await fetch(`${API_URL}/api/geo`, {
      method: 'GET',
      headers: { 'Accept': 'application/json' }
    });
    
    if (!response.ok) return null;
    
    const data = await response.json();
    console.log('☁️ Cloudflare country:', data.country);
    return data.country && data.country !== 'XX' ? data.country : null;
  } catch (error) {
    console.error('Failed to detect country from backend:', error);
    return null;
  }
}

export async function detectCountryByIPFallback(): Promise<string | null> {
  try {
    const response = await fetch('http://ip-api.com/json/', {
      method: 'GET',
    });
    
    if (!response.ok) return null;
    
    const data = await response.json();
    console.log('🌐 IP-API Geolocation:', data.countryCode, data.city);
    return data.countryCode || null;
  } catch (error) {
    console.error('Failed to detect country by IP-API:', error);
    return null;
  }
}

export async function detectCountryByIP(): Promise<string | null> {
  try {
    const response = await fetch('https://ipapi.co/json/', {
      method: 'GET',
      headers: { 'Accept': 'application/json' }
    });
    
    if (!response.ok) return null;
    
    const data = await response.json();
    console.log('🌐 IPAPI.co Geolocation:', data.country_code, data.city);
    return data.country_code || null;
  } catch (error) {
    console.error('Failed to detect country by ipapi.co:', error);
    return null;
  }
}

export async function detectCountry(): Promise<string> {
  if (typeof window === "undefined") return "CH";

  let detectedCountry: string | null = null;

  detectedCountry = await detectCountryByIP();
  if (detectedCountry) {
    console.log('✅ Country detected from ipapi.co:', detectedCountry);
    return detectedCountry;
  }

  detectedCountry = await detectCountryByIPFallback();
  if (detectedCountry) {
    console.log('✅ Country detected from ip-api.com:', detectedCountry);
    return detectedCountry;
  }

  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  console.log('🌍 Detecting country from timezone:', timezone);

  const timezoneMap: Record<string, string> = {
    "Europe/Zurich": "CH", "Europe/Berlin": "DE", "Europe/Vienna": "AT", "Europe/Paris": "FR",
    "Europe/Rome": "IT", "Europe/Madrid": "ES", "Europe/Lisbon": "PT", "Europe/Amsterdam": "NL",
    "Europe/Brussels": "BE", "Europe/Warsaw": "PL", "Europe/Prague": "CZ", "Europe/Stockholm": "SE",
    "Europe/Oslo": "NO", "Europe/Copenhagen": "DK", "Europe/Helsinki": "FI", "Europe/London": "GB",
    "Europe/Dublin": "IE", "Europe/Athens": "GR", "Europe/Budapest": "HU", "Europe/Bucharest": "RO",
    "Europe/Sofia": "BG", "Europe/Vilnius": "LT", "Europe/Riga": "LV", "Europe/Tallinn": "EE",
    "Europe/Ljubljana": "SI", "Europe/Bratislava": "SK", "Europe/Zagreb": "HR", "Europe/Belgrade": "RS",
    "Europe/Kiev": "UA", "Europe/Kyiv": "UA", "Europe/Moscow": "RU", "Europe/Istanbul": "TR",
    "Europe/Minsk": "BY", "Europe/Chisinau": "MD", "Europe/Sarajevo": "BA", "Europe/Podgorica": "ME",
    "Europe/Skopje": "MK", "Europe/Tirane": "AL", "Europe/Vaduz": "LI", "Europe/Luxembourg": "LU",
    "Europe/Monaco": "MC", "Europe/San_Marino": "SM", "Europe/Vatican": "VA", "Europe/Andorra": "AD",
    "America/New_York": "US", "America/Los_Angeles": "US", "America/Chicago": "US", "America/Denver": "US",
    "America/Phoenix": "US", "America/Anchorage": "US", "America/Honolulu": "US", "America/Toronto": "CA",
    "America/Vancouver": "CA", "America/Montreal": "CA", "America/Halifax": "CA", "America/Winnipeg": "CA",
    "America/Mexico_City": "MX", "America/Cancun": "MX", "America/Tijuana": "MX", "America/Monterrey": "MX",
    "America/Sao_Paulo": "BR", "America/Rio_de_Janeiro": "BR", "America/Brasilia": "BR", "America/Fortaleza": "BR",
    "America/Buenos_Aires": "AR", "America/Santiago": "CL", "America/Lima": "PE", "America/Bogota": "CO",
    "America/Caracas": "VE", "America/Panama": "PA", "America/Havana": "CU", "America/Santo_Domingo": "DO",
    "America/Guatemala": "GT", "America/Managua": "NI", "America/San_Jose": "CR", "America/Tegucigalpa": "HN",
    "Asia/Tokyo": "JP", "Asia/Seoul": "KR", "Asia/Shanghai": "CN", "Asia/Beijing": "CN", "Asia/Hong_Kong": "HK",
    "Asia/Taipei": "TW", "Asia/Singapore": "SG", "Asia/Bangkok": "TH", "Asia/Ho_Chi_Minh": "VN",
    "Asia/Jakarta": "ID", "Asia/Manila": "PH", "Asia/Kuala_Lumpur": "MY", "Asia/Kolkata": "IN",
    "Asia/Mumbai": "IN", "Asia/Delhi": "IN", "Asia/Dhaka": "BD", "Asia/Karachi": "PK",
    "Asia/Dubai": "AE", "Asia/Riyadh": "SA", "Asia/Tel_Aviv": "IL", "Asia/Istanbul": "TR",
    "Asia/Ankara": "TR", "Asia/Tehran": "IR", "Asia/Baghdad": "IQ", "Asia/Kuwait": "KW",
    "Australia/Sydney": "AU", "Australia/Melbourne": "AU", "Australia/Brisbane": "AU", "Australia/Perth": "AU",
    "Australia/Adelaide": "AU", "Australia/Darwin": "AU", "Australia/Hobart": "AU", "Pacific/Auckland": "NZ",
    "Pacific/Wellington": "NZ", "Pacific/Fiji": "FJ", "Pacific/Guam": "GU", "Pacific/Honolulu": "US",
    "Africa/Cairo": "EG", "Africa/Johannesburg": "ZA", "Africa/Lagos": "NG", "Africa/Nairobi": "KE",
    "Africa/Casablanca": "MA", "Africa/Algiers": "DZ", "Africa/Tunis": "TN", "Africa/Tripoli": "LY",
    "Africa/Accra": "GH", "Africa/Addis_Ababa": "ET", "Africa/Dar_es_Salaam": "TZ", "Africa/Kampala": "UG",
  };

  if (timezoneMap[timezone]) {
    console.log('✅ Country detected from timezone:', timezoneMap[timezone]);
    return timezoneMap[timezone];
  }

  console.log('⚠️ Timezone not found in map, trying locale...');
  const locale = navigator.language || "en-US";
  console.log('🗣️ Browser locale:', locale);
  const localeMap: Record<string, string> = {
    "de-CH": "CH", "fr-CH": "CH", "it-CH": "CH", "rm-CH": "CH", "de-DE": "DE", "de-AT": "AT",
    "fr-FR": "FR", "it-IT": "IT", "es-ES": "ES", "pt-PT": "PT", "nl-NL": "NL",
    "nl-BE": "BE", "fr-BE": "BE", "pl-PL": "PL", "cs-CZ": "CZ", "sv-SE": "SE",
    "no-NO": "NO", "da-DK": "DK", "fi-FI": "FI", "en-GB": "GB", "en-US": "US",
    "en-CA": "CA", "fr-CA": "CA", "es-MX": "MX", "pt-BR": "BR", "es-AR": "AR",
    "es-CL": "CL", "es-CO": "CO", "es-PE": "PE", "ja-JP": "JP", "ko-KR": "KR",
    "zh-CN": "CN", "zh-HK": "HK", "zh-TW": "TW", "th-TH": "TH", "vi-VN": "VN",
    "id-ID": "ID", "ms-MY": "MY", "en-SG": "SG", "en-PH": "PH", "hi-IN": "IN",
    "en-IN": "IN", "ar-AE": "AE", "ar-SA": "SA", "he-IL": "IL", "tr-TR": "TR",
    "en-AU": "AU", "en-NZ": "NZ", "ar-EG": "EG", "en-ZA": "ZA", "en-NG": "NG",
    "el-GR": "GR", "hu-HU": "HU", "ro-RO": "RO", "bg-BG": "BG", "hr-HR": "HR",
    "sk-SK": "SK", "sl-SI": "SI", "et-EE": "EE", "lv-LV": "LV", "lt-LT": "LT",
    "uk-UA": "UA", "ru-RU": "RU", "sr-RS": "RS", "en-IE": "IE", "sw-KE": "KE",
    "ar-MA": "MA", "ar-DZ": "DZ", "ar-TN": "TN",
  };

  if (localeMap[locale]) {
    console.log('✅ Country detected from locale:', localeMap[locale]);
    return localeMap[locale];
  }

  console.log('❌ Country not found, defaulting to US');
  return "US";
}

export function detectLanguage(): string {
  if (typeof window === "undefined") return "en";

  const locale = navigator.language || "en";
  const langCode = locale.split("-")[0];

  const supportedLanguages = [
    "de", "fr", "it", "en", "es", "pt", "nl", "pl", "cs", "sv", "no", "da", "fi",
    "ja", "ko", "zh", "th", "vi", "id", "ms", "hi", "ar", "he", "tr", "el", "hu",
    "ro", "bg", "hr", "sk", "sl", "et", "lv", "lt", "uk", "ru", "sr", "sw",
    "ka", "az", "be", "bs", "ca", "cy", "eo", "eu", "fa", "fil", "ga", "gl",
    "gu", "ha", "hy", "is", "jv", "ka", "kk", "km", "kn", "ku", "ky", "lo",
    "mk", "ml", "mn", "mr", "mt", "my", "ne", "pa", "ps", "sd", "si", "so",
    "sq", "ta", "te", "tg", "tk", "tl", "ur", "uz", "xh", "yi", "yo", "zu",
    "af", "am", "bn", "ceb", "co", "fy", "gd", "haw", "hmn", "ig", "iw", "jw",
    "kn", "la", "lb", "mg", "mi", "ny", "or", "sm", "sn", "st", "su", "ti"
  ];

  if (supportedLanguages.includes(langCode)) {
    return langCode;
  }

  return "en";
}

export function getCurrencyForCountry(country: string): string {
  const currencyMap: Record<string, string> = {
    CH: "CHF", DE: "EUR", AT: "EUR", FR: "EUR", IT: "EUR", ES: "EUR", PT: "EUR",
    NL: "EUR", BE: "EUR", PL: "PLN", CZ: "CZK", SE: "SEK", NO: "NOK", DK: "DKK",
    FI: "EUR", GB: "GBP", US: "USD", CA: "CAD", MX: "MXN", BR: "BRL", AR: "ARS",
    CL: "CLP", CO: "COP", PE: "PEN", VE: "VES", PA: "PAB", CU: "CUP", JP: "JPY",
    KR: "KRW", CN: "CNY", HK: "HKD", TW: "TWD", SG: "SGD", TH: "THB", VN: "VND",
    ID: "IDR", PH: "PHP", MY: "MYR", IN: "INR", AE: "AED", SA: "SAR", IL: "ILS",
    TR: "TRY", AU: "AUD", NZ: "NZD", EG: "EGP", ZA: "ZAR", NG: "NGN", KE: "KES",
    MA: "MAD", DZ: "DZD", TN: "TND", IE: "EUR", GR: "EUR", HU: "HUF", RO: "RON",
    BG: "BGN", HR: "EUR", SK: "EUR", SI: "EUR", EE: "EUR", LV: "EUR", LT: "EUR",
    UA: "UAH", RU: "RUB", RS: "RSD",
  };

  return currencyMap[country] || "USD";
}