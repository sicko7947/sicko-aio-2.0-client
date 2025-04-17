package constants

// CountryInfo : CountryInfo
type CountryInfo struct {
	Country    string `json:"country,omitempty"`
	MerchGroup string `json:"merchGroup,omitempty"`
	Language   string `json:"language,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Locale     string `json:"locale,omitempty"`
	Domain     string `json:"domain,omitempty"`
}
