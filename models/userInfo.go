package models

// ProfileGroupName : `json:"profileGroupName"`
type ProfileGroupName string

// ProfileName : `json:"profileName"`
type ProfileName string

// AccountGroupName : `json:"accountGroupName"`
type AccountGroupName string

// ProxyGroupName :  `json:"proxyGroupName"`
type ProxyGroupName string

// DiscountGroupName : `json:"discountGroupName"`
type DiscountGroupName string

// GiftCardGroupName : `json:"giftCardGroupName"`
type GiftCardGroupName string

// Profile : profile data
type Profile struct {
	ProfileName string `json:"profileName"`
	Email       string `json:"email"`

	ShippingFirstName   string `json:"shippingFirstName"`
	ShippingLastName    string `json:"shippingLastName"`
	ShippingAddress1    string `json:"shippingAddress1"`
	ShippingAddress2    string `json:"shippingAddress2"`
	ShippingCounty      string `json:"shippingCounty,omitempty"`
	ShippingCity        string `json:"shippingCity"`
	ShippingState       string `json:"shippingState"`
	ShippingProvince    string `json:"shippingProvince,omitempty"`
	ShippingCountry     string `json:"shippingCountry"`
	ShippingZipCode     string `json:"shippingZipCode"`
	ShippingPhoneNumber string `json:"shippingPhoneNumber"`

	BillingFirstName   string `json:"billingFirstName"`
	BillingLastName    string `json:"billingLastName"`
	BillingAddress1    string `json:"billingAddress1"`
	BillingAddress2    string `json:"billingAddress2"`
	BillingCounty      string `json:"billingCounty,omitempty"`
	BillingCity        string `json:"billingCity"`
	BillingState       string `json:"billingState"`
	BillingProvince    string `json:"billingProvince,omitempty"`
	BillingCountry     string `json:"billingCountry"`
	BillingZipCode     string `json:"billingZipCode"`
	BillingPhoneNumber string `json:"billingPhoneNumber"`

	CardHolderName   string `json:"cardHolderName"`
	CardType         string `json:"cardType"`
	CreditCardNumber string `json:"creditCardNumber"`
	ExpireMonth      string `json:"expireMonth"`
	ExpireYear       string `json:"expireYear"`
	CardCVV          string `json:"cardCVV"`
}

// Account : account info
type Account struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	UserID       string `json:"userId,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Status       string `json:"status,omitempty"`
	LastSyncTime string `json:"lastSyncTime,omitempty"`
}

// GiftCard : gift card info
type GiftCard struct {
	Region     string  `json:"region"`
	CardNumber string  `json:"cardNumber"`
	CardPin    string  `json:"cardPin"`
	Balance    float64 `json:"balance,omitempty"`
}

// Settings : general bot settings
type Settings struct {
	Key           string              `json:"key"`
	SizeGroup     map[string][]string `json:"sizeGroup"`
	DiscordWebook map[CATEGORY]string `json:"discordWebhook"`
	SlackWebhok   map[CATEGORY]string `json:"slackWebhook"`
}

// Config : all config data
type Config struct {
	TaskGroups map[TaskGroupID]*TaskGroup                    `json:"TaskGroups"`
	Profiles   map[ProfileGroupName]map[ProfileName]*Profile `json:"Profiles"`
	Accounts   map[AccountGroupName][]*Account               `json:"Accounts"`
	GiftCards  map[GiftCardGroupName][]*GiftCard             `json:"Giftcards"`
	Proxies    map[ProxyGroupName][]string                   `json:"Proxies"`
	Discounts  map[DiscountGroupName][]string                `json:"Discounts"`
	Settings   *Settings                                     `json:"Settings"`
}
