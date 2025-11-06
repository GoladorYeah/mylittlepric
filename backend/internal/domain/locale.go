package domain

import "mylittleprice/internal/constants"

// ═══════════════════════════════════════════════════════════
// LOCALE & REGION TYPES
// ═══════════════════════════════════════════════════════════

// CountryCode represents ISO country code
type CountryCode string

const (
	CountryCH CountryCode = "CH" // Switzerland
	CountryDE CountryCode = "DE" // Germany
	CountryAT CountryCode = "AT" // Austria
	CountryFR CountryCode = "FR" // France
	CountryIT CountryCode = "IT" // Italy
	CountryES CountryCode = "ES" // Spain
	CountryGB CountryCode = "GB" // United Kingdom
	CountryUS CountryCode = "US" // United States
)

// LanguageCode represents ISO language code
type LanguageCode string

const (
	LanguageDE LanguageCode = "de" // German
	LanguageFR LanguageCode = "fr" // French
	LanguageIT LanguageCode = "it" // Italian
	LanguageEN LanguageCode = "en" // English
	LanguageES LanguageCode = "es" // Spanish
)

// Currency represents currency code
type Currency string

const (
	CurrencyCHF Currency = "CHF" // Swiss Franc
	CurrencyEUR Currency = "EUR" // Euro
	CurrencyGBP Currency = "GBP" // British Pound
	CurrencyUSD Currency = "USD" // US Dollar
)

// Locale represents user's regional settings
type Locale struct {
	Country  CountryCode
	Language LanguageCode
	Currency Currency
}

// NewLocale creates a new Locale with defaults
func NewLocale(country, language string) Locale {
	countryCode := CountryCode(country)
	if countryCode == "" {
		countryCode = CountryCode(constants.DefaultCountry)
	}

	langCode := LanguageCode(language)
	if langCode == "" {
		langCode = GetDefaultLanguage(countryCode)
	}

	return Locale{
		Country:  countryCode,
		Language: langCode,
		Currency: GetCurrencyForCountry(countryCode),
	}
}

// String returns string representation
func (l Locale) String() string {
	return string(l.Country) + "_" + string(l.Language)
}

// GetCurrencyForCountry maps country to its currency
func GetCurrencyForCountry(country CountryCode) Currency {
	currencyMap := map[CountryCode]Currency{
		CountryCH: CurrencyCHF,
		CountryDE: CurrencyEUR,
		CountryAT: CurrencyEUR,
		CountryFR: CurrencyEUR,
		CountryIT: CurrencyEUR,
		CountryES: CurrencyEUR,
		CountryGB: CurrencyGBP,
		CountryUS: CurrencyUSD,
	}

	if currency, ok := currencyMap[country]; ok {
		return currency
	}
	return CurrencyEUR
}

// GetDefaultLanguage maps country to its default language
func GetDefaultLanguage(country CountryCode) LanguageCode {
	languageMap := map[CountryCode]LanguageCode{
		CountryCH: LanguageDE,
		CountryDE: LanguageDE,
		CountryAT: LanguageDE,
		CountryFR: LanguageFR,
		CountryIT: LanguageIT,
		CountryES: LanguageES,
		CountryGB: LanguageEN,
		CountryUS: LanguageEN,
	}

	if lang, ok := languageMap[country]; ok {
		return lang
	}
	return LanguageEN
}

// GetLanguageName returns full language name
func GetLanguageName(lang LanguageCode) string {
	nameMap := map[LanguageCode]string{
		LanguageDE: "German",
		LanguageFR: "French",
		LanguageIT: "Italian",
		LanguageEN: "English",
		LanguageES: "Spanish",
	}

	if name, ok := nameMap[lang]; ok {
		return name
	}
	return "English"
}
