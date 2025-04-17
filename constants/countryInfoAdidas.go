package constants

// GetAdidasCountryInfo : Get adidas country info
func GetAdidasCountryInfo(country string) *CountryInfo {

	switch country {
	case "AU":
		return &CountryInfo{
			Domain: "www.adidas.com.au",
		}
	case "CA":
		return &CountryInfo{
			Domain: "www.adidas.ca",
		}
	case "US":
		return &CountryInfo{
			Domain: "www.adidas.com",
		}
	case "MX":
		return &CountryInfo{
			Domain: "www.adidas.mx",
		}
	case "ZA":
		return &CountryInfo{
			Domain: "www.adidas.co.za",
		}
	case "AR":
		return &CountryInfo{
			Domain: "www.adidas.com.ar",
		}
	case "BR":
		return &CountryInfo{
			Domain: "www.adidas.com.br",
		}
	case "CO":
		return &CountryInfo{
			Domain: "www.adidas.co",
		}
	case "PE":
		return &CountryInfo{
			Domain: "www.adidas.pe",
		}
	case "NZ":
		return &CountryInfo{
			Domain: "www.adidas.co.nz",
		}
	case "TH":
		return &CountryInfo{
			Domain: "www.adidas.co.th",
		}
	case "PH":
		return &CountryInfo{
			Domain: "www.adidas.com.ph",
		}
	case "VN":
		return &CountryInfo{
			Domain: "www.adidas.com.vn",
		}
	case "MY":
		return &CountryInfo{
			Domain: "www.adidas.com.my",
		}
	case "SG":
		return &CountryInfo{
			Domain: "www.adidas.com.sg",
		}
	case "PL":
		return &CountryInfo{
			Domain: "www.adidas.pl",
		}
	case "FR":
		return &CountryInfo{
			Domain: "www.adidas.fr",
		}
	case "NL":
		return &CountryInfo{
			Domain: "www.adidas.nl",
		}
	case "UK":
		return &CountryInfo{
			Domain: "www.adidas.co.uk",
		}
	case "DE":
		return &CountryInfo{
			Domain: "www.adidas.de",
		}
	case "IT":
		return &CountryInfo{
			Domain: "www.adidas.it",
		}
	case "GR":
		return &CountryInfo{
			Domain: "www.adidas.gr",
		}
	case "ES":
		return &CountryInfo{
			Domain: "www.adidas.es",
		}
	default:
		return &CountryInfo{
			Domain: "www.adidas.com",
		}
	}
}
