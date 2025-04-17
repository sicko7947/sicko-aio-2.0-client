package constants

// GetLouisVuittonCountryInfo : Get louis vuitton country info
func GetLouisVuittonCountryInfo(country string) *CountryInfo {

	switch country {

	// AMERICAS
	case "US":
		return &CountryInfo{
			Domain:   "us.louisvuitton.com",
			Currency: "USD",
		}
	case "BR":
		return &CountryInfo{
			Domain:   "br.louisvuitton.com",
			Currency: "BRL",
		}
	case "CA":
		return &CountryInfo{
			Domain:   "ca.louisvuitton.com",
			Currency: "CAD",
		}
	// OCEANIA
	case "AU":
		return &CountryInfo{
			Domain:   "au.louisvuitton.com",
			Currency: "AUD",
		}
	case "NZ":
		return &CountryInfo{
			Domain:   "nz.louisvuitton.com",
			Currency: "NZD",
		}

	// ASIA PACIFIC
	case "CN":
		return &CountryInfo{
			Domain:   "www.louisvuitton.cn",
			Currency: "CNY",
		}
	case "HK":
		return &CountryInfo{
			Domain:   "hk.louisvuitton.com/",
			Currency: "HKD",
		}
	case "TW":
		return &CountryInfo{
			Domain:   "tw.louisvuitton.com/",
			Currency: "TWD",
		}
	case "JP":
		return &CountryInfo{
			Domain:   "kr.louisvuitton.com",
			Currency: "JPY",
		}
	case "KR":
		return &CountryInfo{
			Domain:   "kr.louisvuitton.com",
			Currency: "KRW",
		}
	case "TH":
		return &CountryInfo{
			Domain:   "th.louisvuitton.com/",
			Currency: "THB",
		}
	case "MY":
		return &CountryInfo{
			Domain:   "ap.louisvuitton.com/",
			Currency: "MYR",
		}
	case "SG":
		return &CountryInfo{
			Domain:   "ap.louisvuitton.com/",
			Currency: "SGD",
		}

	// EU
	case "UK":
		return &CountryInfo{
			Domain:   "uk.louisvuitton.com",
			Currency: "GBP",
		}
	case "DE":
		return &CountryInfo{
			Domain:   "de.louisvuitton.com",
			Currency: "EUR",
		}
	case "IT":
		return &CountryInfo{
			Domain:   "it.louisvuitton.com",
			Currency: "EUR",
		}
	case "FR":
		return &CountryInfo{
			Domain:   "fr.louisvuitton.com",
			Currency: "EUR",
		}
	case "ES":
		return &CountryInfo{
			Domain:   "es.louisvuitton.com",
			Currency: "EUR",
		}
	case "RU":
		return &CountryInfo{
			Domain:   "ru.louisvuitton.com",
			Currency: "RUB",
		}
	case "DK", "BE", "IE", "LU", "MC", "NL", "AT", "FI", "SE":
		return &CountryInfo{
			Domain:   "en.louisvuitton.com",
			Currency: "EUR",
		}

	default:
		return &CountryInfo{
			Domain: "www.louisvuitton.com",
		}
	}
}
