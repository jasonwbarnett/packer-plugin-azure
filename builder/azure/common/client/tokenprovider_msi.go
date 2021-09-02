package client

import (
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

// for managed identity auth
type msiOAuthTokenProvider struct {
	env      azure.Environment
	clientID string
}

func NewMSIOAuthTokenProvider(env azure.Environment, clientID string) oAuthTokenProvider {
	return &msiOAuthTokenProvider{env: env, clientID: clientID}
}

func (tp *msiOAuthTokenProvider) getServicePrincipalToken() (*adal.ServicePrincipalToken, error) {
	return tp.getServicePrincipalTokenWithResource(tp.env.ResourceManagerEndpoint)
}

func (tp *msiOAuthTokenProvider) getServicePrincipalTokenWithResource(resource string) (*adal.ServicePrincipalToken, error) {
	if tp.clientID != "" {
		return adal.NewServicePrincipalTokenFromMSIWithUserAssignedID("http://169.254.169.254/metadata/identity/oauth2/token", resource, tp.clientID)
	} else {
		return adal.NewServicePrincipalTokenFromMSI("http://169.254.169.254/metadata/identity/oauth2/token", resource)
	}
}
