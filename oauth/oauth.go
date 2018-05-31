package oauth

import (
	"github.com/vkuznecovas/mouthful/config/model"
	"github.com/vkuznecovas/mouthful/oauth/provider"
)

// GetProviders returns a slice of providers formed from the corresponding config section
func GetProviders(providers *[]model.OauthProvider, cbURIBase string) ([]provider.Provider, error) {
	result := make([]provider.Provider, 0)
	if providers == nil {
		return result, nil
	}

	for _, v := range *providers {
		if v.Enabled == false {
			continue
		}

		uri := cbURIBase + v.Name
		p, err := provider.New(v.Name, v.Secret, v.Key, *v.AdminUserIds, uri)
		if err != nil {
			return result, err
		}
		result = append(result, *p)
	}
	return result, nil
}
