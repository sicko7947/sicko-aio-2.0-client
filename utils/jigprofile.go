package utils

import (
	"github.com/sicko7947/sickocommon"
	"sicko-aio-2.0-client/models"
)

func JigProfile(profile *models.Profile) *models.Profile {

	profile.Email = sickocommon.Jig(profile.Email)

	profile.ShippingFirstName = sickocommon.Jig(profile.ShippingFirstName)
	profile.ShippingLastName = sickocommon.Jig(profile.ShippingLastName)
	profile.ShippingAddress1 = sickocommon.Jig(profile.ShippingAddress1)
	profile.ShippingAddress2 = sickocommon.Jig(profile.ShippingAddress2)
	profile.ShippingCounty = sickocommon.Jig(profile.ShippingCounty)
	profile.ShippingCity = sickocommon.Jig(profile.ShippingCity)
	profile.ShippingState = sickocommon.Jig(profile.ShippingState)
	profile.ShippingProvince = sickocommon.Jig(profile.ShippingProvince)
	profile.ShippingCountry = sickocommon.Jig(profile.ShippingCountry)
	profile.ShippingZipCode = sickocommon.Jig(profile.ShippingZipCode)
	profile.ShippingPhoneNumber = sickocommon.Jig(profile.ShippingPhoneNumber)

	profile.CardHolderName = sickocommon.Jig(profile.CardHolderName)

	if profile.BillingAddress1 == "" {
		profile.BillingFirstName = profile.ShippingFirstName
		profile.BillingLastName = profile.ShippingLastName
		profile.BillingAddress1 = profile.ShippingAddress1
		profile.BillingAddress2 = profile.ShippingAddress2
		profile.BillingCounty = profile.ShippingCounty
		profile.BillingCity = profile.ShippingCity
		profile.BillingState = profile.ShippingState
		profile.BillingProvince = profile.ShippingProvince
		profile.BillingCountry = profile.ShippingCountry
		profile.BillingZipCode = profile.ShippingZipCode
		profile.BillingPhoneNumber = profile.ShippingPhoneNumber
		return profile
	}

	profile.BillingFirstName = sickocommon.Jig(profile.BillingFirstName)
	profile.BillingLastName = sickocommon.Jig(profile.BillingLastName)
	profile.BillingAddress1 = sickocommon.Jig(profile.BillingAddress1)
	profile.BillingAddress2 = sickocommon.Jig(profile.BillingAddress2)
	profile.BillingCounty = sickocommon.Jig(profile.BillingCounty)
	profile.BillingCity = sickocommon.Jig(profile.BillingCity)
	profile.BillingState = sickocommon.Jig(profile.BillingState)
	profile.BillingProvince = sickocommon.Jig(profile.BillingProvince)
	profile.BillingCountry = sickocommon.Jig(profile.BillingCountry)
	profile.BillingZipCode = sickocommon.Jig(profile.BillingZipCode)
	profile.BillingPhoneNumber = sickocommon.Jig(profile.BillingPhoneNumber)

	return profile
}
